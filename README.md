# 🚀 BWORK

**BWORK** es un framework moderno escrito en **Go** que permite crear APIs robustas, rápidas y modulares en minutos. Con una estructura clara inspirada en Django y un CLI potente, BWORK está diseñado para escalar sin perder simplicidad.

---

## 🧰 Características

- CLI integrado para generar módulos completos (`model + controller + view`)
- Enrutamiento automático
- Soporte para middlewares
- Estructura limpia y extensible
- Configuración por entorno vía `.env`
- Inspirado en buenas prácticas y Clean Architecture

---

## ⚙️ Instalación del CLI

Instala el CLI globalmente con:

```bash
go install github.com/alun-dra/bwork/cmd/bwork@latest
```

---

## 🛠 Comandos disponibles

```bash
bwork init                         # Inicializa un nuevo proyecto
bwork generate model User          # Crea solo un modelo
bwork generate controller User     # Crea solo un controlador
bwork generate view User           # Crea una vista (endpoint)
bwork generate module User         # Crea model + controller + view
bwork run                          # Ejecuta el servidor
```

---

## 📁 Estructura generada

Al ejecutar `bwork init myproyect`, se crea una estructura como:

```
back/
├── bwork_modules/
│   └── router/         # Núcleo del enrutador
├── controllers/        # Controladores
├── models/             # Estructuras de datos
├── views/              # Endpoints
├── routes.go           # Registro automático de rutas
├── main.go
└── .env                # Variables de entorno
```

---

## 🌱 Ejemplo rápido

```bash
bwork generate module usuario
bwork run
```

Visita: [http://localhost:8081/usuario](http://localhost:8081/usuario)

---

## 🔁 Registro automático

BWORK se encarga de registrar las rutas automáticamente en `routes.go`:

```go
func SetupRoutes(mux *http.ServeMux) {
  // Las rutas se insertan automáticamente aquí 🚀
}
```

---

## 🧪 Variables de entorno

El archivo `.env` contiene la configuración básica para conectarte a tu base de datos MySQL:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=bworkdb
```

---

## 💖 ¿Te gusta el proyecto?

Si este framework te ha ayudado o quieres apoyar su desarrollo, ¡considera apoyar en GitHub Sponsors!

> Próximamente: más plantillas, middlewares listos para producción y documentación interactiva ✨

---

## 📄 Licencia

MIT © [alun-dra](https://github.com/alun-dra)
