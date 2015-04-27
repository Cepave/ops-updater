package http

import (
	"net/http"
)

func configProcRoutes() {

	// 这个文件中主要放置一些调试接口，展示内存状态，/proc/echo/只是个占位的，其他方法可以拷贝这个模板
	http.HandleFunc("/proc/echo/", func(w http.ResponseWriter, r *http.Request) {
		arg := r.URL.Path[len("/proc/echo/"):]
		w.Write([]byte(arg))
	})

}
