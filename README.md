# AtaxxGUI
GUI to play Ataxx against a human or computer

### Features
So far there is no support to play against an engine or different time controls

It supports a notation inspired in chess [fen](https://en.wikipedia.org/wiki/Forsyth%E2%80%93Edwards_Notation) which makes it easier to work with engines. The protocol which will be used to interact with them will be inspired by [uci](http://wbec-ridderkerk.nl/html/UCIProtocol.html)

### Use
Compile the code (this needs the use of the [qt go binding](https://github.com/therecipe/qt) and go itself). You may need to install qt5 in your computer.

$ `./main -path /home/name/go/AtaxxGUI-master`

Default is without time control ("inf"), to play with a certain time control use the `-tc` flag. It is in seconds and what comes after the '+' is the increment

$ `./main -path /home/name/go/AtaxxGUI-master -tc 12+5`

Where the first argument has to be the path to the directory, this is to ensure the sprites are loaded. This also offers the possibility for the end user to manually change the sprites.

It works in linux, support is no guaranteed in other OSs

### TODOS

 * Make it work with engines
 * Have a counter of the number of pieces each side has
 * Clean up code, possibly using more files