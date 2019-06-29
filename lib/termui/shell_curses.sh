#!/usr/bin/env bash

################################################################
function addch
{
    addstr "${1:0:1}"
    return ${?}
}
################################################################
function addstr
{
    [[ "_${1}" != "_" ]] &&
      BUF_SCREEN="${BUF_SCREEN}${1}"
    return ${?}
}
################################################################
function attroff
{
    addstr "${CMD_ATTROFF}"
    return ${?}
}
################################################################
function attron
{
    return 0
}
################################################################
function attrset
{
    addstr "$( ${CMD_ATTRSET} ${1} )"
    return ${?}
}
################################################################
function beep
{
    addstr "${CMD_BEEP}"
    return ${?}
}
################################################################
function chkcols
{
    chkint ${1} ${2} &&
      (( ${2} >= 0 )) &&
      (( ${2} <= ${MAX_COLS} )) &&
      return 0

    ROW_NBR="24"
    COL_NBR="1"

    eval addstr \"${CMD_MOVE}\" &&
      clrtoeol &&
      addstr "${1}: Invalid column number" >&2 &&
      refresh &&
      ${ERROR_PAUSE} &&
      eval addstr \"${CMD_MOVE}\" &&
      clrtoeol &&
      refresh

    return 1
}
################################################################
function chkint
{
    let '${2} + 0' > ${DEV_NULL} 2>&1 &&
      return 0

    ROW_NBR="24"
    COL_NBR="1"

    eval addstr \"${CMD_MOVE}\" &&
      clrtoeol &&
      addstr "${1}: argument not a number" >&2 &&
      refresh &&
      ${ERROR_PAUSE} &&
      eval addstr \"${CMD_MOVE}\" &&
      clrtoeol &&
      refresh

    return 1
}
################################################################
function chklines
{
    chkint ${1} ${2} &&
      (( ${2} >= 0 )) &&
      (( ${2} <= ${MAX_LINES} )) &&
      return 0

    ROW_NBR="24"
    COL_NBR="1"

    eval addstr \"${CMD_MOVE}\" &&
      clrtoeol &&
      addstr "${1}: Invalid line number" >&2 &&
      refresh &&
      ${ERROR_PAUSE} &&
      eval addstr \"${CMD_MOVE}\" &&
      clrtoeol &&
      refresh

    return 1
}
################################################################
function chkparm
{
    [[ "_${2}" = "_" ]] &&
      move 24 1 &&
      clrtoeol &&
      addstr "${1}: Missing parameter" >&2 &&
      refresh &&
      ${ERROR_PAUSE} &&
      move 24 1 &&
      clrtoeol &&
      return 1

    return 0
}
################################################################
function clear
{
    addstr "${CMD_CLEAR}"
    return ${?}
}
################################################################
function clrtobol
{
    addstr "${CMD_CLRTOBOL}"
    return ${?}
}
################################################################
function clrtobot
{
    addstr "${CMD_CLRTOEOD}"
    return ${?}
}
################################################################
function clrtoeol
{
    addstr "${CMD_CLRTOEOL}"
    return ${?}
}
################################################################
function delch
{
    addstr "${CMD_DELCH}"
    return ${?}
}
################################################################
function deleteln
{
    addstr "${CMD_DELETELN}"
    return ${?}
}
################################################################
function endwin
{
    unset MAX_LINES
    unset MAX_COLS
    unset BUF_SCREEN
    return ${?}
}
################################################################
function getch
{
    IFS='' read -r -- TMP_GETCH
    STATUS="${?}"
#     ${CMD_ECHO} "${TMP_GETCH}"
    eval \${CMD_ECHO} ${OPT_ECHO} \"\${TMP_GETCH}\"
    return ${STATUS}
}
################################################################
function getstr
{
    IFS="${IFS_CR}"
    getch
    STATUS="${?}"
    IFS="${IFS_NORM}"
    return ${STATUS}
}
################################################################
function getwd
{
    getch
    return ${?}
}
################################################################
function initscr
{
    PGMNAME="Bourne Shell Curses demo"
    DEV_NULL="/dev/null"
    CMD_TPUT="tput"			# Terminal "put" command

    eval CMD_MOVE=\`echo \"`tput cup`\" \| sed \\\
-e \"s/%p1%d/\\\\\${1}/g\" \\\
-e \"s/%p2%d/\\\\\${2}/g\" \\\
-e \"s/%p1%02d/\\\\\${1}/g\" \\\
-e \"s/%p2%02d/\\\\\${2}/g\" \\\
-e \"s/%p1%03d/\\\\\${1}/g\" \\\
-e \"s/%p2%03d/\\\\\${2}/g\" \\\
-e \"s/%p1%03d/\\\\\${1}/g\" \\\
-e \"s/%d\\\;%dH/\\\\\${1}\\\;\\\\\${2}H/g\" \\\
-e \"s/%p1%c/'\\\\\\\`echo \\\\\\\${1} P | dc\\\\\\\`'/g\" \\\
-e \"s/%p2%c/'\\\\\\\`echo \\\\\\\${2} P | dc\\\\\\\`'/g\" \\\
-e \"s/%p1%\' \'%+%c/'\\\\\\\`echo \\\\\\\${1} 32 + P | dc\\\\\\\`'/g\" \\\
-e \"s/%p2%\' \'%+%c/'\\\\\\\`echo \\\\\\\${2} 32 + P | dc\\\\\\\`'/g\" \\\
-e \"s/%p1%\'@\'%+%c/'\\\\\\\`echo \\\\\\\${1} 100 + P | dc\\\\\\\`'/g\" \\\
-e \"s/%p2%\'@\'%+%c/'\\\\\\\`echo \\\\\\\${2} 100 + P | dc\\\\\\\`'/g\" \\\
-e \"s/%i//g\;s/%n//g\"\`

    CMD_CLEAR="$( ${CMD_TPUT} clear 2>${DEV_NULL} )"	  # Clear display
    CMD_LINES="$( ${CMD_TPUT} lines 2>${DEV_NULL} )"	  # Number of lines on display
    CMD_COLS="$( ${CMD_TPUT} cols 2>${DEV_NULL} )"	  # Number of columns on display
    CMD_CLRTOEOL="$( ${CMD_TPUT} el 2>${DEV_NULL} )"	  # Clear to end of line
    CMD_CLRTOBOL="$( ${CMD_TPUT} el1 2>${DEV_NULL} )"	  # Clear to beginning of line
    CMD_CLRTOEOD="$( ${CMD_TPUT} ed 2>${DEV_NULL} )"	  # Clear to end of display
    CMD_DELCH="$( ${CMD_TPUT} dch1 2>${DEV_NULL} )"	  # Delete current character
    CMD_DELETELN="$( ${CMD_TPUT} dl1 2>${DEV_NULL} )"	  # Delete current line
    CMD_INSCH="$( ${CMD_TPUT} ich1 2>${DEV_NULL} )"	  # Insert 1 character
    CMD_INSERTLN="$( ${CMD_TPUT} il1 2>${DEV_NULL} )"  # Insert 1 Line
    CMD_ATTROFF="$( ${CMD_TPUT} sgr0 2>${DEV_NULL} )"  # All Attributes OFF
    CMD_ATTRSET="${CMD_TPUT}"			  # requires arg ( rev, blink, etc )
    CMD_BEEP="$( ${CMD_TPUT} bel 2>${DEV_NULL} )"	  # ring bell
    CMD_LISTER="cat"
    CMD_SYMLNK="ln -s"
    CMD_ECHO="print"
    CMD_ECHO="echo"
    #OPT_ECHO='-n --'
    OPT_ECHO='-n'
    CMD_MAIL="mail"
    WHOAMI="${LOGNAME}@$( uname -n )"
    WRITER="dfrench@mtxia.com"
    CMD_NOTIFY="\${CMD_ECHO} ${OPT_ECHO} \"\${PGMNAME} - \${WHOAMI} - \$( date )\" | \${CMD_MAIL} \${WRITER}"
    ERROR_PAUSE="sleep 2"

    case "_$( uname -s )" in
        "_Windows_NT") ${DEV_NULL}="NUL";
                       CMD_SYMLNK="cp";;
#              "_Linux") CMD_ECHO="echo -e";;
    esac

    IFS_CR="$'\n'"
    IFS_CR="
"
    IFS_NORM="$' \t\n'"
    IFS_NORM="
"

    MAC_TIME="TIMESTAMP=\`date +\"%y:%m:%d:%H:%M:%S\"\`"
    MAX_LINES=$( ${CMD_TPUT} lines )
    MAX_COLS=$( ${CMD_TPUT} cols )
    BUF_SCREEN=""
    BUF_TOT=""

    return 0
}
################################################################
function insch
{
    addstr "${CMD_INSCH}"
    return ${?}
}
################################################################
function insertln
{
    addstr "${CMD_INSERTLN}"
    return ${?}
}
################################################################
function move
{
#     chklines "${0}" "${1}" \
#     && chkcols "${0}" "${2}" \
#
################################################################
# HEATH-KIT MOVE COMMAND
#    addstr "${1} ${2}"
# VT100 MOVE COMMAND
#    addstr "1};${2}H"
# TPUT MOVE COMMAND
    eval addstr \"${CMD_MOVE}\"
# HP TERMINAL MOVE COMMAND
#   addstr "${1}y${2}C"
################################################################
#  add your move command below this line

    return ${?}
}
################################################################
function mvaddch
{
    move "${1}" "${2}" &&
      addch "${3}"
    return ${?}
}
################################################################
function mvaddstr
{
    move "${1}" "${2}" &&
      addstr "${3}"
    return ${?}
}
################################################################
function mvclrtobol
{
    move "${1}" "${2}" &&
      clrtobol
    return ${?}
}
################################################################
function mvclrtobot
{
    move "${1}" "${2}" &&
      clrtobot
    return ${?}
}
################################################################
function mvclrtoeol
{
    move "${1}" "${2}" &&
      clrtoeol
    return ${?}
}
################################################################
function mvcur
{
    chklines "${0}" "${1}" &&
      chkcols "${0}" "${2}" &&
      eval \"${CMD_MOVE}\"
    return ${?}
}
################################################################
function mvdelch
{
    move "${1}" "${2}" &&
      addstr "${CMD_DELCH}"
    return ${?}
}
################################################################
function mvinsch
{
    move "${1}" "${2}" &&
      addstr "${CMD_INSCH}"
    return ${?}
}
################################################################
function refresh
{
    if [[ "_${1}" != "_" ]]
    then
        eval \${CMD_ECHO} \${OPT_ECHO} \"\${${1}}\"
    else
        ${CMD_ECHO} ${OPT_ECHO} "${BUF_SCREEN}"
        BUF_TOT="${BUF_TOT}${BUF_SCREEN}"
        BUF_SCREEN=""
    fi
    return 0
}
################################################################
function savescr
{
    [[ "_${DEV_NULL}" != "_${1}" ]] &&
      eval ${1}="\"\${BUF_TOT}\""
    BUF_TOT=""
    return ${?}
}
################################################################