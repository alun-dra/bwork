package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func runServer() {
	module := getModuleName() // ← tu función actual, p.ej. "project"
	fmt.Printf("🚀 ejecutando servidor con hot-reload desde %s/main.go...\n", module)

	// arranque inicial
	ctx, cancel := context.WithCancel(context.Background())
	cmd := startChild(ctx, module)

	// bucle de recarga
	changes := make(chan struct{}, 1)
	go watchModule(module, changes)

	for {
		select {
		case <-changes:
			// debounce simple
			time.Sleep(300 * time.Millisecond)
			fmt.Println("🔁 cambio detectado, reiniciando servidor...")

			stopChild(cmd)
			cancel()

			ctx, cancel = context.WithCancel(context.Background())
			cmd = startChild(ctx, module)
		}
	}
}

// -------------------------------------------------------------------
// Lanzar/killear proceso hijo `go run ./<module>`
// -------------------------------------------------------------------

func startChild(ctx context.Context, module string) *exec.Cmd {
	// si existe main.go usamos "./module/main.go", si no, "./module"
	args := []string{"run"}
	if fileExists(filepath.Join(module, "main.go")) {
		args = append(args, "./"+filepath.ToSlash(filepath.Join(module, "main.go")))
	} else {
		args = append(args, "./"+filepath.ToSlash(module))
	}

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

func stopChild(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	// intento suave (Unix)
	if runtime.GOOS != "windows" {
		_ = cmd.Process.Signal(syscall.SIGINT)
		time.Sleep(400 * time.Millisecond)
	}
	_ = cmd.Process.Kill()
}

// -------------------------------------------------------------------
// Watcher por sondeo (stdlib) dentro de la carpeta del módulo
// -------------------------------------------------------------------

func watchModule(module string, out chan<- struct{}) {
	watchDirs := []string{
		".", "controllers", "views", "models", "routes", "config",
	}
	watchExts := []string{".go", ".html", ".tmpl", ".env"}
	ignoreDirs := map[string]bool{
		".git": true, "vendor": true, "node_modules": true, "bwork_modules/router": true,
	}

	prev := map[string]time.Time{}
	root := module

	for {
		cur := map[string]time.Time{}
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			rel, _ := filepath.Rel(root, path)

			// ignorar directorios no deseados
			if d.IsDir() {
				name := d.Name()
				if ignoreDirs[name] {
					return filepath.SkipDir
				}
				// limitar a las carpetas vistas
				if !inWatchedDir(rel, watchDirs) {
					return filepath.SkipDir
				}
				return nil
			}

			// filtrar por extensión
			if !hasExt(rel, watchExts) {
				return nil
			}

			info, e := d.Info()
			if e != nil {
				return nil
			}
			cur[rel] = info.ModTime()

			// disparar si hubo cambios (evita trigger en el primer escaneo)
			if len(prev) > 0 {
				if t, ok := prev[rel]; !ok || info.ModTime().After(t) {
					fmt.Println("📝 archivo modificado:", filepath.ToSlash(filepath.Join(module, rel)))
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

func inWatchedDir(rel string, dirs []string) bool {
	if rel == "." || rel == "" {
		return true
	}
	rel = filepath.ToSlash(rel)
	for _, d := range dirs {
		if d == "." {
			continue
		}
		d = filepath.ToSlash(d)
		if strings.HasPrefix(rel+"/", d+"/") || rel == d {
			return true
		}
	}
	return false
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
