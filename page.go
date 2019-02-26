package main

import (
	"fmt"
)

func page(title, chart01 string, chart02 string, chart03 string, log string, bgColor string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>%s</title>
    %s
    %s
  </head>
  <body>
    <div class="navbar navbar-static-top %s"><div class="container"><h1>%s</h1></div></div>
    <div class="container">
      <div class="row">
        <div class="col-sm">
          %s
          <br/>
          %s
          <br/>
          %s
        </div>
        <div class="col-sm">
          %s
        </div>
      </div>
  </body>
</html>`, title, staticTextTerminalStylesheet, staticTextCdnIncludes, bgColor, title, chart01, chart02, chart03, log)
}
