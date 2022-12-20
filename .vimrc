set fileencodings=utf-8,ucs-bom,gb18030,gbk,gb2312,cp936
set termencoding=utf-8
set encoding=utf-8

set invlist
set mouse=a
set clipboard=autoselect
set clipboard=unnamed
set clipboard=unnamedplus
vnoremap <C-y> "+y
let mapleader='ljj'

nnoremap <Leader>q :quit!<CR>
nnoremap <Leader>wq :wq<CR>

set incsearch
set ignorecase
set nocompatible
set wildmenu
syntax enable
syntax on
filetype indent on
filetype plugin on      
filetype plugin indent on
set foldmethod=syntax
set nofoldenable
set expandtab
set tabstop=4
" INSERT mode
 let &t_SI = "\<Esc>[6 q" . "\<Esc>]12;blue\x7"
" REPLACE mode
let &t_SR = "\<Esc>[3 q" . "\<Esc>]12;black\x7"
 " NORMAL mode
let &t_EI = "\<Esc>[2 q" . "\<Esc>]12;green\x7"
set pastetoggle=<F11>
nmap LE $
set nu                                                                          
set incsearch
let data_dir = has('nvim') ? stdpath('data') . '/site' : '~/.vim'
if empty(glob(data_dir . '/autoload/plug.vim'))
  silent execute '!curl -fLo '.data_dir.'/autoload/plug.vim --create-dirs  https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim'
  autocmd VimEnter * PlugInstall --sync | source $MYVIMRC
endif
if empty(glob('~/.vim/autoload/plug.vim'))
  silent !curl -fLo ~/.vim/autoload/plug.vim --create-dirs
    \ https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
endif

" Run PlugInstall if there are missing plugins
autocmd VimEnter * if len(filter(values(g:plugs), '!isdirectory(v:val.dir)'))
  \| PlugInstall --sync | source $MYVIMRC
\| endif

call plug#begin('~/.vim/plugged')
Plug 'joom/vim-commentary'
Plug 'neoclide/coc.nvim', {'branch': 'release'}
Plug 'neoclide/coc.nvim', {'branch': 'master', 'do': 'yarn install --frozen-lockfile'}
Plug 'SirVer/ultisnips'
Plug 'dylanaraps/fff.vim'
Plug 'ap/vim-buftabline'
Plug 'andrewstuart/vim-kubernetes'
Plug 'dgryski/vim-godef'
Plug 'honza/vim-snippets'
Plug 'majutsushi/tagbar'
Plug 'rhysd/vim-clang-format'
Plug  'pboettch/vim-cmake-syntax'
Plug 'richq/vim-cmake-completion'
Plug 'richq/cmake-lint'

call plug#end()


"""snippets
imap <C-l> <Plug>(coc-snippets-expand)
vmap <C-p> <Plug>(coc-snippets-select)
let g:coc_snippet_next = '<c-p>'
let g:coc_snippet_prev = '<c-n>'
imap <C-p> <Plug>(coc-snippets-expand-jump)

let g:coc_global_extensions =[
            \ 'coc-json',
            \'coc-vimlsp',
            \'coc-snippets',
             \ 'coc-yaml',
            \'coc-go']
" Some servers have issues with backup files, see #649.
set nobackup
set nowritebackup
" delays and poor user experience.
set updatetime=200
" Always show the signcolumn, otherwise it would shift the text each time
" diagnostics appear/become resolved.
set signcolumn=yes
" Use tab for trigger completion with characters ahead and navigate.
" NOTE: Use command ':verbose imap <tab>' to make sure tab is not mapped by
" other plugin before putting this into your config.


 inoremap <silent><expr> <TAB>
       \ coc#pum#visible() ? coc#pum#next(1):
       \ CheckBackspace() ? "\<Tab>" :
       \ coc#refresh()
" inoremap <expr><S-TAB> coc#pum#visible() ? coc#pum#prev(1) : "\<C-h>"

"" Make <CR> to ::waccept selected completion item or notify coc.nvim to format
" <C-g>u breaks current undo, please make your own choice.
inoremap <silent><expr> <CR> coc#pum#visible() ? coc#pum#confirm()
                              \: "\<C-g>u\<CR>\<c-r>=coc#on_enter()\<CR>"

 function! CheckBackspace() abort
   let col = col('.') - 1
   return !col || getline('.')[col - 1]  =~# '\s'
 endfunction

" " Use <c-space> to trigger completion.
" if has('nvim')
"       inoremap <silent><expr> <c-space> coc#refresh()
"   else
"       inoremap <silent><expr> <c-@> coc#refresh()
"     endif

" nmap <silent> [g <Plug>(coc-diagnostic-prev)
" nmap <silent> ]g <Plug>(coc-diagnostic-next)

" " GoTo code navigation.
" nmap <silent> gd <Plug>(coc-definition)
" nmap <silent> gy <Plug>(coc-type-definition)
" nmap <silent> gi <Plug>(coc-implementation)
" nmap <silent> gr <Plug>(coc-references)

" " Use K to show documentation in preview window.
" nnoremap <silent> K :call ShowDocumentation()<CR>
" function! ShowDocumentation()
"   if CocAction('hasProvider', 'hover')
"     call CocActionAsync('doHover')
"   else
"     call feedkeys('K', 'in')
"   endif
" endfunction

" nmap <leader>rn <Plug>(coc-rename)


" """""""" fff
" " Open fff on press of 'f'
nnoremap f :F<CR>
" " Vertical split (NERDtree style).
 let g:fff#split = "50vnew"
" " Open split on the left side (NERDtree style).
 let g:fff#split_direction = "nosplitbelow nosplitright"

" """ 
" """Plug 'ap/vim-buftabline'
" """ 
 set hidden
 nnoremap <C-N> :bnext<CR>
 nnoremap <C-M> :bprev<CR>
" autocmd BufWritePre *.go :silent call CocAction('runCommand', 'editor.action.organizeImport')
" """
 autocmd FileType go nnoremap <buffer> gd :call GodefUnderCursor()<cr>
 autocmd FileType go nnoremap <buffer> <C-]> :call GodefUnderCursor()<cr>
 let g:godef_split=5
 let g:godef_same_file_in_same_window=1 "
set tags=tags
nmap <F8> :TagbarToggle<CR>
let g:tagbar_width=40
set laststatus=2  
 function! CurDir()
           let curdir = substitute(getcwd(), $HOME, "~", "g")
     return curdir
     endfunction
      set statusline=[%n]\ %f%m%r%h\ \|\ \ pwd:\ %{CurDir()}\ \ \|%=\|\%l,%c\ %p%%\ \|\ ascii=%b,hex=%b%{((&fenc==\"\")?\"\":\"\ \|\ \".&fenc)}\
                 \|\ %{$USER}\ @\ %{hostname()} 
      set showcmd 
      set cmdheight=1
      set showmatch 
      set ignorecase
      set hlsearch
    
"" format_better
map <F3> :call Format()<CR>
func! Format()
        exec "w"
        if &filetype == 'c'
        exec "!clang-format -style=WebKit -i %"
        endif
        endfunc
"""
"""
"""
"""
"""
"""
"""
"""
"""
"""
"""
"""
"""


