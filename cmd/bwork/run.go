package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// --- configuración del watcher ---
var (
	watchDirs = []string{
		".", "controllers", "views", "models", "routes", "config",
	}
	watchExts    = []string{".go", ".html", ".tmpl", ".env"}
	ignoreDirs   = []string{".git", "vendor", "node_modules", "bwork_modules/router"}
	pollInterval = 600 * time.Millisecond
	debounceWin  = 300 * time.Millisecond
)

func runServer() {
	target := getModuleName() // tu función actual
	fmt.Printf("🚀 Ejecutando servidor con hot-reload desde %s...\n", target)

	// arranque inicial
	ctx, cancel := context.WithCancel(context.Background())
	cmd := start(ctx, target)

	// canal para reinicios y para signals
	restartCh := make(chan struct{}, 1)
	stopCh := make(chan os.Signal, 1)
	installSignalHandler(stopCh)

	// watcher en goroutine
	go func() {
		if err := watchAndNotify(restartCh); err != nil {
			fmt.Println("⚠️  watcher error:", err)
		}
	}()

	var mu sync.Mutex
	for {
		select {
		case <-restartCh:
			// debounce simple
			time.Sleep(debounceWin)

			mu.Lock()
			// matar proceso actual
			_ = stop(cmd)
			// nuevo contexto/proceso
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			cmd = start(ctx, target)
			mu.Unlock()

		case <-stopCh:
			_ = stop(cmd)
			cancel()
			return
		}
	}
}

// lanza `go run` y conecta stdio
func start(ctx context.Context, target string) *exec.Cmd {
	args := []string{"run", target}
	// si target es carpeta/archivo normaliza
	if fi, err := os.Stat(target); err == nil && fi.IsDir() {
		args = []string{"run", "./" + target}
	}
	if strings.HasSuffix(strings.ToLower(target), ".go") {
		args = []string{"run", target}
	}

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		fmt.Println("❌ Error al iniciar el servidor:", err)
		return cmd
	}

	go func() {
		_ = cmd.Wait() // evita zombies
	}()

	return cmd
}

// intenta terminar suave y si no, mata
func stop(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	// en *nix intenta SIGINT primero
	if sendInterrupt(cmd) == nil {
		// dar tiempo a cerrar
		done := make(chan struct{})
		go func() { cmd.Wait(); close(done) }()
		select {
		case <-done:
			return nil
		case <-time.After(1 * time.Second):
		}
	}
	// fuerza
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
				// limitar a watchDirs
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
			// cambio nuevo o modificado
			if t, ok := prev[path]; !ok || info.ModTime().After(t) {
				// evita spam inicial: solo dispara si prev ya tenía datos
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
	// "." siempre vale
	if path == "." {
		return true
	}
	// path comienza con cualquiera de los watchDirs
	for _, d := range watchDirs {
		if d == "." {
			continue
		}
		// normaliza separadores
		if strings.HasPrefix(filepath.ToSlash(path)+"/", d+"/") {
			return true
		}
	}
	return false
}

// --- señales cross-platform ---
func installSignalHandler(stopCh chan<- os.Signal) {
	// stdlib soporta SIGINT/SIGTERM; en Windows SIGINT llega con Ctrl+C.
	// usamos Notify en main.go si ya lo tienes; si no:
	go func() {
		// bloquea lectura de Stdin para Ctrl+C (fallback simple)
		buf := make([]byte, 1)
		for {
			_, err := os.Stdin.Read(buf)
			if err != nil {
				stopCh <- os.Interrupt
				return
			}
		}
	}()
}

func sendInterrupt(cmd *exec.Cmd) error {
	// en Go stdlib, en Unix se puede Signal(syscall.SIGINT)
	// para evitar syscall en este snippet y mantener stdlib,
	// usamos Process.Kill si no hay soporte. Si quieres SIGINT:
	//   import "syscall"
	//   return cmd.Process.Signal(syscall.SIGINT)
	return fmt.Errorf("interrupt not sent")
}
