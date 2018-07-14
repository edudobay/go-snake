#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <unistd.h>
#include <SDL/SDL.h>
#include <SDL/SDL_image.h>

/******************************************************************************
 * Utility functions
 ******************************************************************************/

#define ASSERT(cond) \
   do if (!(cond)) { \
      fprintf(stderr, "Assertion failed: %s:%d: %s\n", \
         __FILE__, __LINE__, #cond); \
      exit(1); \
   } while (0)

#define allocate(T,n) mallocX((n)*sizeof(T))


void *mallocX(unsigned size)
{
   void *p;
   p = malloc (size);
   if (p == NULL) {
      fprintf (stderr, "Error allocating memory!\n");
      exit (1);
   }
   return p;
}

/******************************************************************************/

enum _Direction {
   DIR_NONE = 0,
   DIR_UP,
   DIR_RIGHT,
   DIR_DOWN,
   DIR_LEFT,
};

int screen_width = 640, screen_height = 480, screen_bpp = 32;
int sprite_width = 11, sprite_height = 11;
int grid_width = 40, grid_height = 30;
Uint32 bg_color = 0, fg_color = 1;
int current_palette = 0;
int psychedelic_mode = 0;
#define NUM_PALETTES 5

#define MIN_LEVEL 1
#define MAX_LEVEL 9
#define MAX_SUBRATE 8
#define RATE_INCREASE_FREQ 50
const int level_rates[] = { 0, 600, 480, 380, 300, 230, 180, 140, 110, 90 };

typedef enum _Direction Direction;
typedef enum _Sprite Sprite;
typedef enum _MoveResult MoveResult;

enum _Sprite {
   SPR_INVALID = -1,
   SPR_NONE = 0,
   SPR_WALL,
   SPR_FOOD,
   SPR_TURN,
   SPR_BODY,
   SPR_HEAD,
   SPR_FOOD2,
   SPR_FOOD3,
   SPR_BORDER_UL,
   SPR_BORDER_UR,
   SPR_BORDER_BR,
   SPR_BORDER_BL,
   SPR_MARKER,
   SPR_COUNT, /* last one */
};

enum _MoveResult {
   MOVE_OK = 0,
   MOVE_GROW,
   MOVE_WALL,
   MOVE_SELF,
};

enum _State {
   STATE_PLAY,
   STATE_PAUSE,
   STATE_WAIT,
   STATE_INPUT,
   STATE_QUIT,
} state = STATE_PLAY;

   //bg_color = SDL_MapRGB(screen->format, 37, 42, 27);
SDL_Color palettes[][3] = {
   {
      {  74,  83,  53 },
      { 180, 191, 155 },
      {  27,  31,  20 },
   },
   {
      {  83,  53,  74 },
      { 191, 155, 180 },
      {  31,  20,  27 },
   },
   {
      { 115,  38,  38 },
      { 220, 149, 149 },
      {  41,  13,  13 },
   },
   {
      { 115, 113,  38 },
      { 220, 218, 149 },
      {  41,  41,  13 },
   },
   {
      {  38, 115, 115 },
      { 149, 220, 220 },
      {  13,  41,  41 },
   }};

typedef struct _PendingMovement PendingMovement;

struct _PendingMovement {
   Direction direction;
   PendingMovement *next;
};

SDL_Surface *screen = NULL;
SDL_Surface *sprites;
SDL_Rect clips[SPR_COUNT];
SDL_Rect board_clip;
SDL_Rect score_clip;

typedef struct _Font Font;
struct _Font {
   int height;
   SDL_Surface *surface;
   SDL_Rect clips[256];
};
Font *main_font;

typedef enum _TextAnchor TextAnchor;
enum _TextAnchor {
   ANCHOR_LEFT,
   ANCHOR_CENTER,
   ANCHOR_RIGHT,
};

char *board = NULL;
char *map = NULL;
int *board_links = NULL;
Direction cur_direction;
int cur_head_pos;
int cur_size;
int cur_rate;
int cur_subrate;
int frames;
int cur_level;
int move_count;
int points;

int food_count = 0;
int food2_count = 0;
int food3_count = 0;
int has_food = 0;
int has_food2 = 0;
int has_food3 = 0;
int food2_pos;
int food2_start;
int food3_pos;
int food3_start;

#define FOOD2_DURATION 50
#define FOOD3_DURATION 25

void (*wait_callback)(SDL_Event *);

typedef struct _MenuItem MenuItem;
typedef struct _MenuGroup MenuGroup;

struct _MenuItem {
   int x, y;
   const char *text;
   MenuItem *top, *bottom;
   MenuItem *left, *right;
   void (*activate)();
};

struct _MenuGroup {
   //MenuItem *first;
   MenuItem *items;
   int count;
   MenuItem *active;
};

MenuGroup *active_menu = NULL;


PendingMovement *movement_queue = NULL;

void enqueue_movement(Direction direction);
Direction dequeue_movement();
void free_movement_queue(PendingMovement *head);


Uint32 ticks;
void timer_start();
Uint32 timer_get_ticks();

/******************************************************************************
 * Timer functions
 ******************************************************************************/

void timer_start()
{
   ticks = SDL_GetTicks();
}

Uint32 timer_get_ticks()
{
   return SDL_GetTicks() - ticks;
}

/******************************************************************************/

void dump(SDL_Surface *surface, const char *filename)
{
   FILE *out = fopen(filename, "wb");
   SDL_LockSurface(surface);
   //fwrite(surface->pixels, surface->w * surface->h * surface->format->BytesPerPixel, 1, out);
   printf("pitch : %d, w*bpp : %d\n", surface->pitch, surface->w * surface->format->BytesPerPixel);
   fwrite(surface->pixels, surface->h * surface->pitch, 1, out);
   SDL_UnlockSurface(surface);
   fclose(out);
}

/******************************************************************************/

void set_palette(SDL_Surface *surface, int index, int text)
{

   SDL_Color *palette = palettes[index];
   int res;

   if (text) {
      SDL_Color aux;
      palette = malloc(sizeof(palettes[index]));
      memcpy(palette, palettes[index], sizeof(palettes[index]));

      aux = palette[0];
      palette[0] = palette[1];
      palette[1] = aux;
   }
   else {
      bg_color = SDL_MapRGB(screen->format, palette[0].r, palette[0].g, palette[0].b);
      fg_color = SDL_MapRGB(screen->format, palette[1].r, palette[1].g, palette[1].b);
   }

   res = SDL_SetPalette(surface, SDL_LOGPAL | SDL_PHYSPAL, palette, 0, 3);
   if (res == 0) {
      printf("SetPalette -> %d\n", res);
      return;
   }

   if (text)
      free(palette);
}

void load_palette(int index)
{
   set_palette(sprites, index, 0);
   if (main_font)
      set_palette(main_font->surface, index, 1);

   current_palette = index;
}

void video_setup()
{
   screen = SDL_SetVideoMode(screen_width, screen_height, screen_bpp, SDL_HWSURFACE);
   if (screen == NULL) {
      fprintf(stderr, "Unable to set video mode: %s\n", SDL_GetError());
      exit(1);
   }

   SDL_WM_SetCaption("Snake", NULL);
}

void init_sdl()
{
   int img_flags = IMG_INIT_PNG;

   if (SDL_Init(SDL_INIT_AUDIO | SDL_INIT_VIDEO) < 0) {
      fprintf(stderr, "Unable to initialize SDL: %s\n", SDL_GetError());
      exit(1);
   }

   if ((IMG_Init(img_flags) & img_flags) != img_flags) {
      fprintf(stderr, "Unable to initialize image format support: %s\n", IMG_GetError());
      exit(1);
   }

   atexit(SDL_Quit);
   atexit(IMG_Quit);

}

void draw_sprite(int clip, int x, int y)
{
   SDL_Rect dest;

   ASSERT(clip >= 0 && clip < SPR_COUNT);

   dest.x = x;
   dest.y = y;
   dest.w = sprite_width;
   dest.h = sprite_height;

   //if (clip == SPR_NONE)
      //SDL_FillRect(screen, &dest, bg_color);
   //else
   if (clip != SPR_NONE && clip != SPR_INVALID)
      SDL_BlitSurface(sprites, &clips[clip], screen, &dest);
}

void optimize_image(SDL_Surface **image)
{
   SDL_Surface *new;

   ASSERT(image != NULL);
   ASSERT(*image != NULL);

   new = SDL_DisplayFormat(*image);
   if (new == NULL) {
      fprintf(stderr, "SDL_DisplayFormat: conversion failed\n");
      exit(1);
   }

   SDL_FreeSurface(*image);
   *image = new;
}

void copy_as_hwsurface(SDL_Surface **image)
{
   SDL_Surface *old, *new;

   ASSERT(image != NULL);
   ASSERT(*image != NULL);

   old = *image;
   new = SDL_CreateRGBSurface(SDL_HWSURFACE, old->w, old->h, 8, 0, 0, 0, 0);
   if (new == NULL) {
      fprintf(stderr, "SDL_CreateRGBSurface: error\n");
      exit(1);
   }

   memcpy(new->format->palette->colors, old->format->palette->colors,
      old->format->palette->ncolors*sizeof(SDL_Color));
   new->format->palette->ncolors = old->format->palette->ncolors;

   SDL_LockSurface(old);
   SDL_LockSurface(new);
   memcpy(new->pixels, old->pixels, old->h * old->pitch);
   SDL_UnlockSurface(new);
   SDL_UnlockSurface(old);

   SDL_FreeSurface(old);
   *image = new;
}

void load_sprites()
{
   int i, x;

   sprites = IMG_Load("data/sprites.png");
   if (sprites == NULL) {
      fprintf(stderr, "Unable to load sprites: IMG_Load: %s\n", IMG_GetError());
      exit(1);
   }
   //optimize_image(&sprites);
   //copy_as_hwsurface(&sprites);

   /* skip sprite 0 (which means 'none') */
   for (i = 1, x = 0; i < SPR_COUNT; i++, x += sprite_width) {
      clips[i].x = x;
      clips[i].y = 0;
      clips[i].w = sprite_width;
      clips[i].h = sprite_height;
   }

   //set_palette(sprites, current_palette, 0);

   //dump(sprites, "sprites.dump");
}

void free_sprites()
{
   SDL_FreeSurface(sprites);
}

int step_for_direction(Direction direction)
{
   switch (cur_direction) {
      case DIR_UP:    return -grid_width;
      case DIR_DOWN:  return grid_width;
      case DIR_LEFT:  return -1;
      case DIR_RIGHT: return 1;
      default:
         fprintf(stderr, "Invalid direction value caught!\n");
   }

   return 0;
}

MoveResult move()
{
   int x, new_pos, size = grid_width * grid_height;
   int tail, last;
   Direction direction;

   direction = dequeue_movement();
   if (direction != DIR_NONE)
      cur_direction = direction;

   new_pos = cur_head_pos + step_for_direction(cur_direction);
   
   x = cur_head_pos % grid_width;
   if (x == 0 && cur_direction == DIR_LEFT)
      new_pos = cur_head_pos + grid_width - 1;
   else if (x == grid_width - 1 && cur_direction == DIR_RIGHT)
      new_pos = cur_head_pos - grid_width + 1;
   else if (new_pos >= size)
      new_pos = new_pos - size;
   else if (new_pos < 0)
      new_pos = new_pos + size;

   /* crash into wall */
   if (board[new_pos] == SPR_WALL) {
      printf("Crash!\n");
      return MOVE_WALL;
   }
   /* grow snake */
   else if (board[new_pos] == SPR_FOOD) {
      cur_size++;
      board_links[new_pos] = cur_head_pos;
      board[cur_head_pos] = SPR_BODY;
      board[new_pos] = SPR_HEAD;
      cur_head_pos = new_pos;
      has_food--;
      food_count++;
      food2_count++;
      food3_count++;
      points += cur_level;
      return MOVE_GROW;
   }
   else if (board[new_pos] == SPR_FOOD2) {
      cur_size += 2;
      board_links[new_pos] = cur_head_pos;
      board[cur_head_pos] = SPR_BODY;
      board[new_pos] = SPR_HEAD;
      cur_head_pos = new_pos;
      has_food2 = 0;
      points += cur_level * (FOOD2_DURATION - (frames - food2_start) + 1) / 2;
      return MOVE_GROW;
   }
   else if (board[new_pos] == SPR_FOOD3) {
      points += cur_level * (FOOD3_DURATION - (frames - food3_start) + 1);
      has_food3 = 0;
   }

   tail = cur_head_pos;
   last = -1;
   while (board_links[tail] >= 0) {
      last = tail;
      tail = board_links[tail];
   }

   /* self-collision */
   if (board[new_pos] == SPR_BODY && new_pos != tail) {
      printf("Self-collision\n");
      return MOVE_SELF;
   }

   board_links[last] = -1;
   board_links[new_pos] = cur_head_pos;

   board[cur_head_pos] = SPR_BODY;
   board[tail] = SPR_NONE;
   cur_head_pos = new_pos;
   board[cur_head_pos] = SPR_HEAD;

   return MOVE_OK;
}

void rotate(Direction direction)
{
   if (direction == cur_direction)
      return;
   /* opposite directions */
   if (abs(direction - cur_direction) == 2)
      return;

   enqueue_movement (direction);

   //cur_direction = direction;
   //move();
   //timer_start();
}

int rand_int(int min, int max)
{
   double p = ((double) rand()) / RAND_MAX;
   return min + (int)floor(p*(max - min + 1));
}

int rand_pos()
{
   int pos, size = grid_width * grid_height;
   do pos = rand_int(0, size-1);
      while (board[pos] != SPR_NONE);
   return pos;
}

void put_food()
{
   int pos = rand_pos();
   board[pos] = SPR_FOOD;
   has_food = 1;
}

void put_food2()
{
   int pos = rand_pos();
   board[pos] = SPR_FOOD2;
   has_food2 = 1;
   food2_start = frames;
   food2_pos = pos;
}

void put_food3()
{
   int pos = rand_pos();
   board[pos] = SPR_FOOD3;
   has_food3 = 1;
   food3_start = frames;
   food3_pos = pos;
}

void remove_food2()
{
   board[food2_pos] = SPR_NONE;
   has_food2 = 0;
}

void remove_food3()
{
   board[food3_pos] = SPR_NONE;
   has_food3 = 0;
}

void free_font(Font *font)
{
   SDL_FreeSurface(font->surface);
   free(font);
}

int is_empty(const char *buf)
{
   const char *p;
   for (p = buf; *p; p++)
      if (!isspace(*p))
         return 0;
   return 1;
}

Font *load_font(const char *def_file)
{
   unsigned char ch, lastch = 0;
   int x, lastx, height;
   FILE *file;
   char *font_file;
   Font *font;
   char buf[256];


   file = fopen(def_file, "r");
   if (file == NULL) {
      fprintf(stderr, "Error: Unable to load font '%s'\n", def_file);
      return NULL;
   }

   if (fscanf(file, "height %d", &height) < 1 || height <= 0) {
      fprintf(stderr, "Error: Invalid font format\n");
      fclose(file);
      return NULL;
   }

   font_file = strdup(def_file);
   strcpy(font_file + strlen(font_file) - 3, "png");

   font = allocate(Font, 1);
   font->height = height;
   font->surface = IMG_Load(font_file);
   free(font_file);

   if (font->surface == NULL) {
      fprintf(stderr, "Unable to load font: %s\n", IMG_GetError());
      fclose(file);
      return NULL;
   }

   //set_palette(font->surface, current_palette, 1);

   memset(font->clips, 0, 256*sizeof(SDL_Rect));

   while (fgets(buf, 256, file) != NULL) {
      if (sscanf(buf, "%d \"%c\"", &x, &ch) < 2) {
         if (is_empty(buf))
            continue;
         fprintf(stderr, "Error: Invalid character definition: %s\n", buf);
         free_font(font);
         fclose(file);
         return NULL;
      }

      if (lastch)
         font->clips[lastch].w = x - lastx;

      font->clips[ch].x = x;
      font->clips[ch].y = 0;
      font->clips[ch].h = height;

      lastch = ch;
      lastx = x;
   }

   if (lastch)
      font->clips[lastch].w = font->surface->w - lastx;

   fclose(file);

   return font;
}

int render_text(const char *text, int x, int y, Font *font, TextAnchor anchor)
{
   unsigned char ch;
   int i, w = 0;
   SDL_Rect dest;

   dest.x = x;
   dest.y = y;
   dest.h = font->height;

   if (anchor != ANCHOR_LEFT) {
      for (i = 0; text[i]; i++) {
         ch = (unsigned char) text[i];

         /* skip non-existant chars */
         if (font->clips[ch].h == 0)
            continue;

         w += font->clips[ch].w;
      }
      if (anchor == ANCHOR_RIGHT)
         dest.x -= w;
      else // ANCHOR_CENTER
         dest.x -= w/2;
   }

   for (i = 0; text[i]; i++) {
      ch = (unsigned char) text[i];

      /* skip empty chars */
      if (font->clips[ch].h == 0)
         continue;

      dest.w = font->clips[ch].w;
      SDL_BlitSurface(font->surface, &font->clips[ch], screen, &dest);

      dest.x += dest.w;
   }

   if (anchor == ANCHOR_LEFT)
      return dest.x - x;
   return w;
}

int load_map(const char *filename)
{
   int i = 0, ch;
   unsigned width, height, size;
   FILE *file;

   file = fopen(filename, "r");
   if (file == NULL) {
      fprintf(stderr, "Error: Unable to load map '%s'\n", filename);
      return 0;
   }

   if (fscanf(file, "%u %u", &width, &height) < 2) {
      fprintf(stderr, "Error: Invalid map format\n");
      return 0;
   }

   size = width * height;
   if (size == 0) {
      fprintf(stderr, "Error: Invalid map format\n");
      return 0;
   }

   grid_width = width;
   grid_height = height;

   if (map != NULL) {
      free(map);
   }

   map = allocate(char, size);

   while (i < size) {
      ch = fgetc(file);
      if (ch == EOF) {
         fprintf(stderr, "Error: Invalid map format (premature end of file)\n");
         free(map);
         map = NULL;
         return 0;
      }

      if (isspace(ch))
         continue;

      switch (ch) {
         case '.': map[i] = SPR_NONE; break;
         case '#': map[i] = SPR_WALL; break;
         case 'x': map[i] = SPR_INVALID; break;
         default:
            fprintf(stderr, "Error: Invalid character 0x%x in map\n", (unsigned)ch);
            free(map);
            map = NULL;
            return 0;
      }

      i++;
   }

   return 1;
}

void draw_window(int x, int y, int w, int h)
{
   SDL_Rect rect = { .x = x, .y = y, .w = w, .h = h };

   SDL_FillRect(screen, &rect, fg_color);

   draw_sprite(SPR_BORDER_UL, x, y);
   draw_sprite(SPR_BORDER_UR, x + w - 11, y);
   draw_sprite(SPR_BORDER_BR, x + w - 11, y + h - 11);
   draw_sprite(SPR_BORDER_BL, x, y + h - 11);
}

void draw_score()
{
   char buf[32];

   SDL_FillRect(screen, &score_clip, fg_color);

   sprintf(buf, "LEVEL %d", cur_level);
   if (cur_subrate > 0)
      sprintf(buf + strlen(buf), "+%d", cur_subrate);

   render_text(buf, score_clip.x + 10, score_clip.y + 5, main_font, ANCHOR_LEFT);


   sprintf(buf, "POINTS %5d", points);
   render_text(buf, score_clip.x + score_clip.w - 10, score_clip.y + 5, main_font, ANCHOR_RIGHT);

   if (has_food2) {
      int remaining = FOOD2_DURATION - (frames - food2_start);
      sprintf(buf, "%d", remaining);
      render_text(buf, score_clip.x + score_clip.w/2, score_clip.y + 28, main_font, ANCHOR_RIGHT);
   }
   else if (has_food3) {
      int remaining = FOOD3_DURATION - (frames - food3_start);
      sprintf(buf, "%d", remaining);
      render_text(buf, score_clip.x + score_clip.w/2, score_clip.y + 28, main_font, ANCHOR_RIGHT);
   }

}

void init_board(int head_pos)
{
   int i, size, pos, xc, yc, step;

   if (board != NULL)
      free(board);
   if (board_links != NULL)
      free(board_links);

   size = grid_width * grid_height;

   board_clip.x = 5;
   board_clip.y = 5;
   board_clip.w = sprite_width * grid_width;
   board_clip.h = sprite_height * grid_height;

   score_clip.x = 5;
   score_clip.y = board_clip.h + 10;
   score_clip.w = board_clip.w;
   score_clip.h = 56;

   screen_width = board_clip.w + 10;
   screen_height = board_clip.h + score_clip.h + 15;

   board = allocate(char, size);

   if (map != NULL)
      memcpy(board, map, size);
   else
      memset(board, SPR_NONE, size);

   board_links = allocate(int, size);
   for (i = 0; i < size; i++)
      board_links[i] = -1;

   xc = grid_width / 2;
   yc = grid_height / 2;

   if (head_pos < 0)
      head_pos = xc + grid_width*yc;

   cur_head_pos = head_pos;

   board[head_pos] = SPR_HEAD;
   step = -step_for_direction(cur_direction);

   for (i = 1, pos = head_pos; i < cur_size; i++) {
      board_links[pos] = pos + step;
      board[pos += step] = SPR_BODY;
   }
}

void draw_board()
{
   int i, j, I = 0, x, y;

   SDL_FillRect(screen, &board_clip, bg_color);

   for (i = 0; i < grid_height; i++) {
      I = i*grid_width;
      for (j = 0; j < grid_width; j++) {
         x = 5 + sprite_width*j;
         y = 5 + sprite_height*i;
         draw_sprite(board[I+j], x, y);
      }
   }
}

void free_movement_queue(PendingMovement *head)
{
   PendingMovement *node = head, *next;

   while (node != NULL) {
      next = node->next;
      free(node);
      node = next;
   }
}

void enqueue_movement(Direction direction)
{
   PendingMovement *node = movement_queue, *new;

   ASSERT (movement_queue != NULL);

   new = allocate(PendingMovement, 1);
   new->direction = direction;
   new->next = NULL;

   while (node->next != NULL)
      node = node->next;

   node->next = new;
}

Direction dequeue_movement()
{
   PendingMovement *next;

   ASSERT (movement_queue != NULL);

   next = movement_queue->next;
   if (next == NULL)
      return DIR_NONE;

   movement_queue->next = next->next;

   if (next->direction == cur_direction)
      return DIR_NONE;
   /* opposite directions */
   if (abs(next->direction - cur_direction) == 2)
      return DIR_NONE;

   move_count++;

   return next->direction;
}

void set_level(int level)
{
   if (!(level >= MIN_LEVEL && level <= MAX_LEVEL)) {
      fprintf(stderr, "Warning: Invalid level number\n");
      return;
   }

   cur_rate = level_rates[level];
   cur_subrate = 0;
   cur_level = level;
}

void increase_subrate()
{
   cur_subrate++;
   cur_rate /= 1.05;
}

void init_game(int level)
{
   cur_direction = DIR_DOWN;
   cur_size = 3;
   frames = 0;
   food_count = 0;
   has_food = 0;
   has_food2 = 0;
   has_food3 = 0;
   points = 0;
   move_count = 0;
   set_level(level);

   timer_start();

   if (movement_queue == NULL) {
      movement_queue = allocate(PendingMovement, 1);
      movement_queue->next = NULL;
      movement_queue->direction = DIR_NONE;
   }
   else {
      free_movement_queue(movement_queue->next);
   }
}

void print_usage(int argc, char *argv[])
{
   fprintf(stderr, "Usage: %s [-l level] [-m mapfile]\n",
      argv[0]);
}

void switch_color()
{
   if (current_palette == NUM_PALETTES-1)
      load_palette(0);
   else
      load_palette(current_palette+1);
}

void restart_game()
{
   init_game(cur_level);
   init_board(-1);

   draw_score();
   draw_board();
   SDL_Flip(screen);

   state = STATE_PLAY;
}

void quit_game()
{
   state = STATE_QUIT;
}

void draw_menu(MenuGroup *menu)
{
   int i, x, y;
   
   for (i = 0; i < menu->count; i++) {
      x = menu->items[i].x;
      y = menu->items[i].y;

      if (menu->active == &menu->items[i])
         render_text("\002", x, y, main_font, ANCHOR_LEFT);
      else
         render_text(" ", x, y, main_font, ANCHOR_LEFT);

      render_text(menu->items[i].text, x + 17, y, main_font, ANCHOR_LEFT);
   }
}

int translate_key(SDL_keysym sym)
{
   int key = sym.sym;

   if (sym.mod & KMOD_SHIFT || sym.mod & KMOD_CAPS)
   {
      if (key >= SDLK_a && key <= SDLK_z)
         key = key - 32;
      if (sym.sym & KMOD_SHIFT) {
         switch (key) {
            case '1': key = '!'; break;
            case '2': key = '@'; break;
            case '3': key = '#'; break;
            case '4': key = '$'; break;
            case '5': key = '%'; break;
            case '6': key = '^'; break;
            case '7': key = '&'; break;
            case '8': key = '*'; break;
            case '9': key = '('; break;
            case '0': key = ')'; break;
            case '-': key = '_'; break;
            case '=': key = '+'; break;
            case '[': key = '{'; break;
            case '{': key = '}'; break;
            case ';': key = '}'; break;
            case ',': key = '<'; break;
            case '.': key = '>'; break;
            case '/': key = '?'; break;
            case '\\': key = '|'; break;
         }   
      }
   }

   if (sym.mod & KMOD_NUM) {
      if (key >= SDLK_KP0 && key <= SDLK_KP9)
         key = key - 208; /* - 256 + 48 */
      else if (key == SDLK_KP_PERIOD)
         key = '.';
   }
   else {
      switch (key) {
         case SDLK_KP0: key = SDLK_INSERT; break;
         case SDLK_KP1: key = SDLK_END; break;
         case SDLK_KP2: key = SDLK_DOWN; break;
         case SDLK_KP3: key = SDLK_PAGEDOWN; break;
         case SDLK_KP4: key = SDLK_LEFT; break;
         case SDLK_KP6: key = SDLK_RIGHT; break;
         case SDLK_KP7: key = SDLK_HOME; break;
         case SDLK_KP8: key = SDLK_UP; break;
         case SDLK_KP9: key = SDLK_PAGEUP; break;
         case SDLK_KP_PERIOD: key = SDLK_DELETE; break;
      }
   }

   switch (key) {
      case SDLK_KP_DIVIDE: key = '/'; break;
      case SDLK_KP_MULTIPLY: key = '*'; break;
      case SDLK_KP_MINUS: key = '-'; break;
      case SDLK_KP_PLUS: key = '+'; break;
   }

   return key;
}

void menu_destroy(MenuGroup *menu)
{
   if (active_menu == menu)
      active_menu = NULL;

   free(menu->items);
   free(menu);
}

int menu_callback(SDL_Event *event)
{
   MenuItem *active = active_menu->active;
   if (!active)
      return 0;

   if (event->type == SDL_KEYDOWN) {
      switch (event->key.keysym.sym) {
         case SDLK_UP:
            if (active->top)
               active_menu->active = active->top;
            break;
         case SDLK_DOWN:
            if (active->bottom)
               active_menu->active = active->bottom;
            break;
         case SDLK_LEFT:
            if (active->left)
               active_menu->active = active->left;
            break;
         case SDLK_RIGHT:
            if (active->right)
               active_menu->active = active->right;
            break;
         case SDLK_RETURN:
         case SDLK_KP_ENTER:
            if (active->activate) {
               active->activate ();
               //menu_destroy(active_menu);
            }
            break;
         default:
            return 0;
      }
   }

   return 1;
}

void game_over_callback(SDL_Event *event)
{
   if (menu_callback(event)) {
      draw_menu(active_menu);
      SDL_Flip(screen);
      return;
   }

   if (event->type == SDL_KEYDOWN) {
      switch (event->key.keysym.sym) {
         case SDLK_n:
            restart_game();
            break;
         case SDLK_q:
            quit_game();
            break;
         default:
            break;
      }
   }
}

void game_over_display()
{
   int win_x = 71, win_y = 126;
   int win_w = screen_width - 2*win_x, win_h = screen_height - 2*win_y;

   MenuItem *items = allocate(MenuItem, 2);
   MenuGroup *menu = allocate(MenuGroup, 1);

   memset(items, 0, sizeof(MenuItem)*2);

   items[0].text = "NEW GAME";
   items[0].x = win_x + 27;
   items[0].y = win_y + 56;
   items[0].activate = restart_game;
   items[0].bottom = &items[1];
   items[0].top = &items[1];

   items[1].text = "QUIT";
   items[1].x = win_x + 27;
   items[1].y = win_y + 80;
   items[1].activate = quit_game;
   items[1].top = &items[0];
   items[1].bottom = &items[0];

   //.first = &items[0],
   menu->items = items;
   menu->count = 2;
   menu->active = &items[0];

   draw_window(win_x, win_y, win_w, win_h);
   render_text("GAME OVER", win_x + win_w/2, win_y + 12, main_font, ANCHOR_CENTER);

   //render_text("NEW GAME", win_x + 28, win_y + 56, main_font, ANCHOR_LEFT);
   //render_text("QUIT", win_x + 28, win_y + 80, main_font, ANCHOR_LEFT);
   //render_text("\002", win_x + 12, win_y + 56, main_font, ANCHOR_LEFT);

   active_menu = menu;

   draw_menu(menu);
   SDL_Flip(screen);
}

void level_select()
{
   int index = active_menu->active - active_menu->items;

   cur_level = index + 1;
   restart_game();
}

void select_level_callback(SDL_Event *event)
{
   if (menu_callback(event)) {
      draw_menu(active_menu);
      SDL_Flip(screen);
      return;
   }

   if (event->type == SDL_KEYDOWN) {
      int key = translate_key(event->key.keysym);
      if (key >= '1' && key <= '9') {
         active_menu->active = &active_menu->items[key - '1'];
         draw_menu(active_menu);
         SDL_Flip(screen);
      }
      else if (key == SDLK_ESCAPE) {
         menu_destroy(active_menu);
         state = STATE_PLAY;
      }
   }
}

void pause_display()
{
   int win_x = 71, win_h = 48;
   int win_w = screen_width - 2*win_x,
       win_y = (screen_height - win_h) / 2;

   draw_window(win_x, win_y, win_w, win_h);
   render_text("GAME PAUSED", win_x + win_w/2, win_y + 12, main_font, ANCHOR_CENTER);
   SDL_Flip(screen);
}

void select_level_display()
{
   int win_x = 71, win_y = 126;
   int win_w = screen_width - 2*win_x, win_h = screen_height - 2*win_y;
   int i, j;
   
   MenuItem *items = allocate(MenuItem, 9);
   MenuGroup *menu = allocate(MenuGroup, 1);
   const char *menu_texts[] = {"1", "2", "3", "4", "5", "6", "7", "8", "9"};

   draw_window(win_x, win_y, win_w, win_h);
   render_text("SELECT LEVEL", win_x + win_w/2, win_y + 12, main_font, ANCHOR_CENTER);

   memset(items, 0, sizeof(MenuItem)*9);

   for (i = 0; i < 9; i++) {
      items[i].text = menu_texts[i];
      items[i].x = win_x + (win_w - 120)/2;
      items[i].y = win_y + 56;
      items[i].activate = level_select;
      items[i].bottom = &items[(i >= 6) ? (i - 6) : (i + 3)];
      items[i].top = &items[(i < 3) ? (i + 6) : (i - 3)];
      items[i].right = &items[(i == 8) ? 0 : (i + 1)];
      items[i].left = &items[(i == 0) ? 8 : (i - 1)];
   }

   for (i = 0; i < 3; i++) {
      for (j = 0; j < 3; j++) {
         items[3*i + j].x += j * 40;
         items[3*i + j].y += i * 30;
      }
   }

   menu->items = items;
   menu->count = 9;
   menu->active = &items[cur_level-1];

   active_menu = menu;

   draw_menu(menu);
   SDL_Flip(screen);
}

int main(int argc, char *argv[])
{
   int yes_food2 = 0, yes_food3 = 0;
   int init_level = 7;
   int opt;
   char *map_name = NULL;

   while ((opt = getopt(argc, argv, "l:m:")) != -1) {
      switch (opt) {
         case 'l':
            init_level = atoi(optarg);
            break;
         case 'm':
            map_name = malloc(strlen(optarg) + 9);
            strcpy(map_name, "data/");
            strcat(map_name, optarg);
            strcat(map_name, ".map");
            break;
         default:
            print_usage(argc, argv);
            return 1;
      }
   }

   init_sdl();

   init_game(init_level);

   if (map_name) {
      load_map(map_name);
      free(map_name);
   }

   init_board(-1);
   video_setup();

   load_sprites();
   main_font = load_font("data/font1.def");
   load_palette(0);

   draw_score();
   SDL_Flip(screen);

   while (state != STATE_QUIT) {
      SDL_Event event;
      while (SDL_PollEvent(&event)) {

         if (event.type == SDL_QUIT) {
            state = STATE_QUIT;
            break;
         }

         if (state == STATE_PLAY) {
            if (event.type == SDL_KEYDOWN) {
               switch (event.key.keysym.sym) {
                  case SDLK_UP:
                     enqueue_movement (DIR_UP);
                     break;
                  case SDLK_DOWN:
                     enqueue_movement (DIR_DOWN);
                     break;
                  case SDLK_LEFT:
                     enqueue_movement (DIR_LEFT);
                     break;
                  case SDLK_RIGHT:
                     enqueue_movement (DIR_RIGHT);
                     break;
                  case SDLK_q:
                     state = STATE_QUIT;
                     break;
                  case SDLK_p:
                  case SDLK_PAUSE:
                     state = STATE_PAUSE;
                     pause_display();
                     break;
                  /*case SDLK_d:
                     dump(screen, "screen.dump");
                     break;*/
                  case SDLK_f:
                     yes_food2 = 1;
                     break;
                  case SDLK_g:
                     yes_food3 = 1;
                     break;
                  case SDLK_F2:
                     switch_color();
                     break;
                  case SDLK_F3:
                     state = STATE_WAIT;
                     wait_callback = select_level_callback;
                     select_level_display();
                     break;
                  case SDLK_F4:
                     psychedelic_mode = !psychedelic_mode;
                  default:
                     //printf("Pressed key %s\n", SDL_GetKeyName(event.key.keysym.sym));
                     break;
               }
            }
         }

         else if (state == STATE_PAUSE) {
            if (event.type == SDL_KEYDOWN) {
               switch (event.key.keysym.sym) {
                  case SDLK_q:
                     state = STATE_QUIT;
                     break;
                  case SDLK_p:
                  case SDLK_PAUSE:
                     state = STATE_PLAY;
                     break;
                  case SDLK_F2:
                     switch_color();
                     break;
                  default:
                     break;
               }
            }
         }

         else if (state == STATE_WAIT) {
            if (wait_callback)
               wait_callback(&event);
         }

      }

      if (state == STATE_QUIT)
         break;

      if (state == STATE_PLAY) {
         draw_board();
         SDL_Flip(screen);

         if (timer_get_ticks() > cur_rate) {
            MoveResult result = move();

            if (result == MOVE_OK) {
               timer_start();
               frames++;
            }
            else if (result == MOVE_GROW) {
               timer_start();
               frames++;
            }
            else {
               state = STATE_WAIT;
               wait_callback = game_over_callback;
               game_over_display();
               continue;
            }

            if (move_count == RATE_INCREASE_FREQ && cur_subrate < MAX_SUBRATE) {
               move_count = 0;
               increase_subrate();
            }

            if (psychedelic_mode && frames % 2 == 0)
               switch_color();

            draw_score();

         }

         if (!has_food)
            put_food();

         if (!has_food3 && (yes_food3 || food3_count == 15)) {
            put_food3();
            yes_food3 = 0;
            food2_count = 0;
            food3_count = 0;
         }
         if (!has_food2 && (yes_food2 || food2_count == 5)) {
            put_food2();
            yes_food2 = 0;
            food2_count = 0;
         }

         if (has_food2 && frames - food2_start >= FOOD2_DURATION) {
            remove_food2();
         }
         if (has_food3 && frames - food3_start >= FOOD3_DURATION) {
            remove_food3();
         }
      }
   }

   free_font(main_font);
   free_sprites();

   return 0;
}

