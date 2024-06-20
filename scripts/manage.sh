#!/usr/bin/env bash

# Run by: './scripts/manage.sh' from project root

Version=1.0.0

# >>>>>>>>>>>>>>>>>>>>>>>> functions >>>>>>>>>>>>>>>>>>>>>>>>

function install () {

  echo "Building binaries..."
  go build -race -o secretserviced cmd/app/secretserviced/main.go
  go build -race -o secretservice cmd/app/secretservice/main.go
  echo "Copying binaries to /usr/bin"
  # Alternatively: ~/.local/bin
  sudo cp secretserviced /usr/bin
  sudo cp secretservice /usr/bin
  echo "Creating systemd UNIT file at /etc/systemd/user"
  # Alternatively: ~/.config/systemd/user/
  rm secretserviced
  rm secretservice

  echo

  if [ -f "/etc/systemd/user/secretserviced.service" ]; then
    read -rep "systemd UNIT file exist at '/etc/systemd/user/secretserviced.service'.
Delete it [YOU MAY LOOSE YOUR MASTERPASSWORD!!!] (y/n)? " -i "n" answer

    if [[ "${answer}" == "y" ]]; then
      echo "deleting UNIT file"
      sudo rm -f /etc/systemd/user/secretserviced.service
    fi
    
  fi

  echo

  password=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 32 ; echo '')

  cat << EOF | sudo tee "/etc/systemd/user/secretserviced.service" >/dev/null
[Unit]
Description=Service to keep secrets of applications
Documentation=https://github.com/yousefvand/secret-service

[Install]
WantedBy=default.target

[Service]
Type=simple
RestartSec=30
Restart=always
Environment="MASTERPASSWORD=$password"
WorkingDirectory=/usr/bin/
ExecStart=/usr/bin/secretserviced

EOF

# echo "enabling service..."
# systemctl enable --now --user secretserviced.service

echo "$(tput setaf 2)""Done!""$(tput sgr0)"

}

function uninstall () {
  
  # echo "disabling service..."
  # systemctl disable --now --user secretserviced.service
  echo "deleting binaries"
  sudo rm /usr/bin/secretserviced
  sudo rm /usr/bin/secretservice

  if [ -f "/etc/systemd/user/secretserviced.service" ]; then
    read -rep "systemd UNIT file exist at '/etc/systemd/user/secretserviced.service'.
Delete it [YOU MAY LOOSE YOUR MASTERPASSWORD!!!] (y/n)? " -i "n" answer

    if [[ "${answer}" == "y" ]]; then
      echo "deleting UNIT file"
      sudo rm -f /etc/systemd/user/secretserviced.service
    fi
    
  fi

  echo "$(tput setaf 2)""Done!""$(tput sgr0)"
  
}

function help () {
  echo "
  secretserviced installer/uninstaller
  version: $(tput setaf 6)${Version}$(tput sgr0)
  By: Remisa Yousefvand

  flags:

  -i | --install   | install and enable and start service
  ---+-------------+---------------------------------------
  -u | --uninstall | uninstall and disable and stop service
  ---+-------------+---------------------------------------
  -h | --help      | show usage message
"
}

# Usage: options=("one" "two" "three"); inputChoice "Choose:" 1 "${options[@]}"; choice=$?; echo "${options[$choice]}"
function inputChoice() {
  echo "${1}"; shift
  echo "$(tput dim)""-Change option: [up/down], Select: [ENTER]""$(tput sgr0)"
  local selected="${1}"; shift

  ESC=$(echo -e "\033")
  cursor_blink_on()  { tput cnorm; }
  cursor_blink_off() { tput civis; }
  cursor_to()        { tput cup $(($1-1)); }
  print_option()     { echo "$(tput sgr0)" "$1" "$(tput sgr0)"; }
  print_selected()   { echo "$(tput rev)" "$1" "$(tput sgr0)"; }
  get_cursor_row()   { IFS=';' read -rsdR -p $'\E[6n' ROW COL; echo "${ROW#*[}"; }
  key_input()        { read -rs -n3 key 2>/dev/null >&2; [[ $key = ${ESC}[A ]] && echo up; [[ $key = $ESC[B ]] && echo down; [[ $key = "" ]] && echo enter; }

  for opt; do echo; done

  local lastrow
  lastrow=$(get_cursor_row)
  local startrow=$((lastrow - $#))
  trap "cursor_blink_on; echo; echo; exit" 2
  cursor_blink_off

  : selected:=0

  while true; do
    local idx=0
    for opt; do
      cursor_to $((startrow + idx))
      if [ ${idx} -eq "${selected}" ]; then
        print_selected "${opt}"
      else
        print_option "${opt}"
      fi
      ((idx++))
    done

    case $(key_input) in
      enter) break;;
      up)    ((selected--)); [ "${selected}" -lt 0 ] && selected=$(($# - 1));;
      down)  ((selected++)); [ "${selected}" -ge $# ] && selected=0;;
    esac
  done

  cursor_to "${lastrow}"
  cursor_blink_on
  echo

  return "${selected}"
}


# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> argument parsing >>>>>>>>>>>>>>>>>>>>>>>>

POSITIONAL=()
while (( $# > 0 )); do
  case "${1}" in
    -i|--install)
    install
    shift
    ;;
    -u|--uninstall)
    uninstall
    shift
    ;;
    -h|--help)
    help
    shift
    ;;
    *) # unknown flag/switch
    POSITIONAL+=("${1}")
    shift
    ;;
  esac
done

set -- "${POSITIONAL[@]}" # restore positional params

# <<<<<<<<<<<<<<<<<<<<<<<< argument parsing <<<<<<<<<<<<<<<<<<<<<<<<

# Entry point

options=("install" "uninstall" "quit")
inputChoice "Choose operation:" 0 "${options[@]}"; choice=$?


case "${options[$choice]}" in
  install)
    install
  ;;
  uninstall)
    uninstall
  ;;
  *)
    exit 0
  ;;
esac
