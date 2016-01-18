package feedloggr2

const CSS_BODY = `
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

article {
	border: 1px solid #DDD;
	background-color: #FFF;
	border-radius: 3px;
	margin-bottom: 20px;
}

article > h1 {
	background-color: #EEE;
	border-bottom: 1px solid #DDD;
	padding: 10px 15px;
	margin: 0;
	font-size: 16px;
	font-weight: bold;
	text-align: center;
}

article > ol, p {
	margin: 15px;
}

article > ol > li > a:hover {
	color: #000;
}

article > ol > li > a:visited {
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

const HTML_BODY = `
<!doctype html>
<html>
	<head>
		<meta charset="utf-8">
		<title>
			{{format .Date}} | News
		</title>
		<link href="/style.css" rel="stylesheet" type="text/css">
	</head>

	<body>
		<header>
			<nav>
				<a href="{{date_link 24 .Date}}">&lt;</a>
				<a href="index.html">Latest</a>
				<a href="{{date_link -24 .Date}}">&gt;</a>
			</nav>
		</header>

		<section id="content">
			{{range .Feeds}}
			<article>
				<h1>
					<a href="{{.URL}}" rel="nofollow">{{.Title}}</a>
				</h1>
				{{if .Error}}
				<p>Error while updating feed:<br />{{.Error}}</p>
				{{else}}
				<ol>
					{{range .Items}}
					<li>
						<a href="{{.URL}}" rel="nofollow">
							{{.Title}}
						</a>
					</li>
					{{end}}
				</ol>
				{{end}}
			</article>
			{{else}}
			<p class="center">Sorry, no news for today!</p>
			{{end}}
		</section>

		<footer>
			Generated with Feedloggr2
		</footer>
	</body>
</html>
`
