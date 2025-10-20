package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func runServer() {
	module := getModuleName()
	fmt.Printf("🚀 ejecutando servidor con hot-reload desde %s/main.go...\n", module)

	// canal de cambios
	changes := make(chan struct{}, 1)
	go watch(module, changes)

	// bucle infinito: cada vuelta arranca el server y se queda
	// hasta que llegue un cambio; luego termina y vuelve a arrancar.
	for {
		runOnce(module, changes)
	}
}

// ---- una “ronda” del servidor: arranca, espera cambio, corta ----

func runOnce(module string, changes <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // ✅ este cancel SIEMPRE se ejecuta (sin warnings)

	cmd := startChild(ctx, module)

	// esperamos el primer cambio que llegue
	<-changes
	time.Sleep(250 * time.Millisecond) // debounce pequeño
	fmt.Println("🔁 cambio detectado, reiniciando servidor...")

	_ = stopChild(cmd)
}

// ---- helpers de proceso ----

// lanza `go run ./<module>` (siempre el paquete, no un archivo suelto)
func startChild(ctx context.Context, module string) *exec.Cmd {
	args := []string{"run", "./" + filepath.ToSlash(module)}

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		fmt.Println("❌ error al iniciar el servidor:", err)
		return cmd
	}
	go func() { _ = cmd.Wait() }()
	return cmd
}

func stopChild(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	// Matar proceso actual
	_ = cmd.Process.Kill()

	// 🔁 Esperar a que el puerto se libere (máx 3 segundos)
	for i := 0; i < 30; i++ {
		if !portInUse(8081) { // usa tu puerto
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("⚠️ puerto 8081 aún ocupado después de matar el proceso")
}

// portInUse comprueba si el puerto está ocupado
func portInUse(port int) bool {
	addr := fmt.Sprintf("localhost:%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return true // está en uso
	}
	_ = l.Close()
	return false
}

// ---- watcher por sondeo (stdlib) ----

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
			rel, _ := filepath.Rel(root, path)

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

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
