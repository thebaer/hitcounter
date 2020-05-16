# hitcounter

Hitcounter is a simple web application for counting page views, written in Go. It's built for privacy-friendly environments and personal sites like [the author's](https://baer.works), when you just want to know that someone in the world saw your website.

The application compiles to a single binary that you can run anywhere. It stores everything in memory, then persists data to a single JSON file upon shutdown. No database or external dependencies are needed.

## Getting Started

Build the application:

```
go get github.com/thebaer/hitcounter/cmd/hitcounter
``` 

Then run it:

```
hitcounter
```

Woohoo!

## Embedding in your site

Embed the following code for invisible counting:

```html
<img src="http://localhost:6767/hit.gif?p=PAGE/GOES/HERE" style="border:0;" alt="" />
```

Replace `http://localhost:6767` with the domain where you're running this application.

Replace `PAGE/GOES/HERE` with the current page's path. It should be different for each page.

### Hit counter

To show your visitors how many page views you've received, embed the following code, replacing certain parts as outlined in the instructions above.

```html
<div id="count">Viewed ??? times</div>

<noscript><img src="http://localhost:6767/hit.gif?p=PAGE/GOES/HERE" style="border:0;" alt="" /></noscript>
<script>
    var xh = new XMLHttpRequest();
    xh.onreadystatechange = function() {
        if (xh.readyState == XMLHttpRequest.DONE && xh.status == 200) {
            document.getElementById('count').innerHTML = "Viewed " + xh.responseText + " times";
        }
    };
    xh.open("GET", "http://localhost:6767/hit?p=PAGE/GOES/HERE", true);
    xh.send();
</script>
```

## API

All endpoints respond to `GET` requests.

### `/hit?p=KEY`

This endpoint counts a new view for the given `KEY` and then returns the total count as plain text.

### `/hit.gif?p=KEY`

This endpoint counts a new view for the given `KEY`, serving an invisible pixel.

### `/hits?p=KEY`

This endpoint returns the total count for the given `KEY` as plain text **without** recording a new view.