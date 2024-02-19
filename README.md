## Structure

`url-shortener` - main application  
`sso` - auth for application  
`protos` - generated proto files  


#### For development

create file `go.work` with content

```go
go 1.21.7

use (
	protos
	sso
	url-shortener
)
```
