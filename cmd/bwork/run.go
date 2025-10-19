package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

// --- configuración del watcher ---
var (
	watchDirs    = []string{".", "controllers", "views", "models", "routes", "config"}
	watchExts    = []string{".go", ".html", ".tmpl", ".env"}
	ignoreDirs   = []string{".git", "vendor", "node_modules", "bwork_modules/router"}
	pollInterval = 600 * time.Millisecond
	debounceWin  = 300 * time.Millisecond
)

func runServer() {
	target := detectEntrypoint() // <- NUEVO
	fmt.Printf("🚀 Ejecutando servidor con hot-reload usando: %s\n", strings.Join(target, " "))

	// arranque inicial
	ctx, cancel := context.WithCancel(context.Background())
	cmd := start(ctx, target)

	// signals
	restartCh := make(chan struct{}, 1)
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGTERM)

	// watcher
	go func() {
		if err := watchAndNotify(restartCh); err != nil {
			fmt.Println("⚠️  watcher error:", err)
		}
	}()

	var mu sync.Mutex
	for {
		select {
		case <-restartCh:
			time.Sleep(debounceWin)
			mu.Lock()
			_ = stop(cmd)
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			cmd = start(ctx, target)
			mu.Unlock()

		case <-quitCh:
			_ = stop(cmd)
			cancel()
			return
		}
	}
}

// Detecta la mejor forma de ejecutar el entrypoint.
// Preferencias: "go run ./main.go" si existe; si no, "go run ."
func detectEntrypoint() []string {
	if _, err := os.Stat("main.go"); err == nil {
		return []string{"run", "./main.go"}
	}
	// soporta proyectos con cmd/app/main.go o app/main.go si quisieras:
	candidates := []string{"cmd/app/main.go", "app/main.go"}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return []string{"run", "./" + c}
		}
	}
	// por defecto, módulo actual
	return []string{"run", "."}
}

// lanza `go run` y conecta stdio
func start(ctx context.Context, args []string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		fmt.Println("❌ Error al iniciar el servidor:", err)
		return cmd
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			fmt.Println("💥 Proceso del servidor terminó con error:", err)
		} else {
			fmt.Println("🔚 Proceso del servidor terminó.")
		}
	}()

	return cmd
}

// intenta terminar suave y si no, mata
func stop(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	// Unix: SIGINT; Windows: Kill directo
	if runtime.GOOS != "windows" {
		_ = cmd.Process.Signal(syscall.SIGINT)
		time.Sleep(400 * time.Millisecond)
	}
	return cmd.Process.Kill()
}

// --- watcher por sondeo (stdlib puro) ---
func watchAndNotify(restart chan<- struct{}) error {
	prev := make(map[string]time.Time)

	for {
		cur := make(map[string]time.Time)
		err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if shouldIgnoreDir(d.Name()) {
					return filepath.SkipDir
				}
				if !isWithinWatchDirs(path) {
					return nil
				}
				return nil
			}
			if !hasWatchedExt(path) {
				return nil
			}
			info, e := d.Info()
			if e != nil {
				return nil
			}
			cur[path] = info.ModTime()
			if t, ok := prev[path]; !ok || info.ModTime().After(t) {
				if len(prev) > 0 {
					fmt.Println("🔁 cambio detectado en:", path)
					select {
					case restart <- struct{}{}:
					default:
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("watch walk error:", err)
		}
		prev = cur
		time.Sleep(pollInterval)
	}
}

func hasWatchedExt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range watchExts {
		if ext == e {
			return true
		}
	}
	return false
}
func shouldIgnoreDir(name string) bool {
	for _, ig := range ignoreDirs {
		if name == ig {
			return true
		}
	}
	return false
}
func isWithinWatchDirs(path string) bool {
	if path == "." {
		return true
	}
	for _, d := range watchDirs {
		if d == "." {
			continue
		}
		if strings.HasPrefix(filepath.ToSlash(path)+"/", d+"/") {
			return true
		}
	}
	return false
}
