## {{.Version.Os}}{{.Version.Version}}路由修改
{{range .Diff}}
> {{.}}
> 
{{end}}
[详情点击此处](http://{{.Host}}:{{.Port}}/nav/{{.Version.Os}}/{{.Version.App}}/{{.Version.Version}})