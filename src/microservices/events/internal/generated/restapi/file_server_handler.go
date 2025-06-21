// This file is safe to edit. Once it exists it will not be overwritten
package restapi

// to add FileServer feature into server edit setupGlobalMiddleware function in restapi package
// simple example
// func setupGlobalMiddleware(handler http.Handler) http.Handler {
//	return addFileServer(handler)
// }

import (
	"log/slog"
	"net/http"
	"path/filepath"
)

// AssetsPath - каталог файлов файл сервера
const AssetsPath = "/assets"

// addFileServer - мидлварная функция. перехватывает все запросы на сервер.
// для запросов с url AssetsPath отдает файлы
// для остальных запросов передает управление дальше в сервер
func addFileServer(next http.Handler) http.Handler {
	lenAssetsPath := len(AssetsPath)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (len(r.URL.Path) > lenAssetsPath) && (r.URL.Path[0:lenAssetsPath] == AssetsPath) {
			h := http.StripPrefix(AssetsPath+"/", http.FileServer(http.Dir("."+AssetsPath)))

			ext := filepath.Ext(r.URL.Path)
			var contentype string

			switch ext {
			case ".css":
				contentype = "text/css"
			case ".js":
				contentype = "text/javascript"
			case ".png":
				contentype = "image/png"
			case ".ico":
				contentype = "image/ico"
			case ".wasm":
				contentype = "application/wasm"
			case ".map":
				contentype = "application/json"
			case ".jpg":
				contentype = "image/jpeg"
			default:
				contentype = "text/html"
			}

			slog.Debug("Restapi: call addFileServer", "r.URL.Path", r.URL.Path, "ext", ext)

			w.Header().Set("Content-Type", contentype)
			h.ServeHTTP(w, r)

			return
		}

		next.ServeHTTP(w, r)
	})
}
