CC = gcc
CFLAGS = -Wall -g
LDFLAGS = -lSDL -lSDL_image -lm

CROSS_CC = i486-mingw32-gcc
CROSS_INCDIR = ${HOME}/win32/include
CROSS_LIBDIR = ${HOME}/win32/lib
CROSS_BINDIR = ${HOME}/win32/bin
CROSS_CFLAGS = -I${CROSS_INCDIR} ${CFLAGS}
CROSS_LDFLAGS = -L${CROSS_LIBDIR} -lmingw32 -lSDLmain -lSDL -lSDL_image -mwindows

all: snake

pack:
	zip -9r snake-win32.zip snake-win32
snake: snake.o
#	$(CC) -o $@ $(LDFLAGS)

snake.o: snake.c

snake.exe: snake.win32.o
	$(CROSS_CC) -o $@ $^ $(CROSS_LDFLAGS)

snake.win32.o: snake.c
	$(CROSS_CC) -c -o $@ $(CROSS_CFLAGS) $^
