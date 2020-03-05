# vscode support for decksh

Place the contents of the vscode directory into ```<home>/.vscode/extensions/ajstarks.decksh-1.0.0```

The support includes syntax coloring by adding:

```
    "editor.tokenColorCustomizations": {
        "textMateRules": [
            {
                "scope": "keyword.other.command.decksh",
                "settings": {
                    "foreground": "#AA0000"
                }
            },
	...
	}
```

to your settings.

Pair matching works for:

* deck/edeck
* slide/eslide
* list/elist
* nlist/elist
* clist/elist
* blist/elist

If you type the keywords ``li``, ``text``, ``ctext`` ``etext`` followed by space a ``""` will be
automatically inserted.