package feedloggr

const tmplCSS string = `
body {
	margin: 20px 40px;
	line-height: 1.5;
	font-size: 16px;
	font-family: monospace;
	background-color: #FFF;
	color: #444;
}
a, a:hover, a:visited {
	text-decoration:none;
	color: #444;
}
a:hover {
	color: #000;
}

nav {
	text-align: center;
	margin-bottom: 20px;
}
nav > a {
	border: 1px solid #DDD;
	border-radius: 3px;
	padding: 6px 12px;
	font-weight: bold;
}
nav > a:hover {
	background-color: #E6E6E6;
	border-color: #ADADAD;
}

section {
	border: 1px solid #DDD;
	background-color: #FFF;
	border-radius: 3px;
	margin-bottom: 20px;
}
section > h1 {
	background-color: #EEE;
	border-bottom: 1px solid #DDD;
	padding: 10px 15px;
	margin: 0;
	font-size: 16px;
	font-weight: bold;
	text-align: center;
}
section > ol, p {
	margin: 15px;
}
section > ol > li > a:visited {
	color: #AAA;
}

footer {
	text-align: center;
	font-size: 12px;
}
.center {
	text-align: center;
}
`

const tmplPage string = `
<!doctype html>
<html>
<head>
	<meta charset="utf-8">
	<title>{{.CurrentDate}} | Feedloggr</title>
	<link href="./style.css" rel="stylesheet" type="text/css">
</head>
<body>
	<nav>
		<a href="{{.NextDate}}.html">&lt;</a>
		<a href="index.html">Latest</a>
		<a href="{{.PrevDate}}.html">&gt;</a>
	</nav>
	{{range .Feeds}}
	<section>
		<h1><a href="{{.URL}}" rel="nofollow">{{.Title}}</a></h1>
		{{if .Error}}
		<p>Error while updating feed:<br />{{.Error}}</p>
		{{else}}
		<ol>{{range .Items}}
			<li><a href="{{.URL}}" rel="nofollow">{{.Title}}</a></li>
		{{end}}</ol>
		{{end}}
	</section>
	{{else}}
	<p class="center">Sorry, no news for today!</p>
	{{end}}
	<footer>
		Generated with <a href="https://github.com/lmas/feedloggr">Feedloggr</a>
	</footer>
</body>
</html>
`
