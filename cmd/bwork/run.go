package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// -------------------- ENTRYPOINT --------------------

func runServer() {
	module := getModuleName()
	port := detectPort(module)
	fmt.Printf("🚀 ejecutando servidor con hot-reload desde %s (puerto %d)...\n", module, port)

	changes := make(chan struct{}, 1)
	go watch(module, changes)

	for {
		ctx, cancel := context.WithCancel(context.Background())
		cmd := startChild(ctx, module)

		<-changes
		time.Sleep(250 * time.Millisecond)
		fmt.Println("🔁 cambio detectado, reiniciando servidor...")

		_ = stopChild(cmd, port)
		cancel()
	}
}

// -------------------- START / STOP --------------------

func startChild(ctx context.Context, module string) *exec.Cmd {
	args := []string{"run", "./" + filepath.ToSlash(module)}
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin

	// 🔎 buscar go.mod y ejecutar desde el root del módulo
	if root, err := findModuleRoot(); err == nil {
		cmd.Dir = root
	} else {
		fmt.Println("⚠️ no se encontró go.mod, ejecutando desde el directorio actual")
	}

	// ✅ forzar módulos ON (evita que Linux busque paquetes en stdlib)
	env := os.Environ()
	env = append(env, "GO111MODULE=on", "GOWORK=off")
	cmd.Env = env

	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		fmt.Println("❌ error al iniciar el servidor:", err)
		return cmd
	}
	go func() { _ = cmd.Wait() }()
	return cmd
}

func stopChild(cmd *exec.Cmd, port int) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	killProcessTree(cmd.Process.Pid)

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if !portInUse(port) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("⚠️ puerto %d aún ocupado tras reinicio", port)
}

// -------------------- MODULE DETECTION --------------------

// busca go.mod subiendo desde el cwd y devuelve su directorio
func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("no se encontró go.mod en esta ruta ni superiores")
}

// -------------------- NETWORK HELPERS --------------------

func portInUse(port int) bool {
	l, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		return true
	}
	_ = l.Close()
	return false
}

func detectPort(module string) int {
	if v := os.Getenv("BWORK_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			return p
		}
	}
	if b, err := os.ReadFile(filepath.Join(module, ".env")); err == nil {
		for _, ln := range strings.Split(string(b), "\n") {
			ln = strings.TrimSpace(ln)
			if ln == "" || strings.HasPrefix(ln, "#") {
				continue
			}
			if strings.HasPrefix(ln, "PORT=") || strings.HasPrefix(ln, "APP_PORT=") {
				kv := strings.SplitN(ln, "=", 2)
				if len(kv) == 2 {
					if p, err := strconv.Atoi(strings.TrimSpace(kv[1])); err == nil {
						return p
					}
				}
			}
		}
	}
	return 8081
}

// -------------------- WATCHER --------------------

func watch(module string, out chan<- struct{}) {
	root := module
	exts := []string{".go", ".html", ".tmpl", ".env"}
	ignore := map[string]bool{".git": true, "vendor": true, "node_modules": true}

	prev := map[string]time.Time{}
	for {
		cur := map[string]time.Time{}
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if ignore[d.Name()] {
					return filepath.SkipDir
				}
				return nil
			}
			if !hasExt(path, exts) {
				return nil
			}
			info, e := d.Info()
			if e != nil {
				return nil
			}
			rel, _ := filepath.Rel(root, path)
			cur[rel] = info.ModTime()

			if len(prev) > 0 {
				if t, ok := prev[rel]; !ok || info.ModTime().After(t) {
					fmt.Println("📝 archivo modificado:", filepath.ToSlash(path))
					select {
					case out <- struct{}{}:
					default:
					}
				}
			}
			return nil
		})
		prev = cur
		time.Sleep(600 * time.Millisecond)
	}
}

func hasExt(path string, exts []string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range exts {
		if ext == e {
			return true
		}
	}
	return false
}
