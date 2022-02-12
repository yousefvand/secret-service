#!/usr/bin/env bash

# Run by: './scripts/publish-aur.sh' from project root

version="v0.1.0"

tempAURDirectory="tmp-aur"

if [ -d "${tempAURDirectory}" ]; then
  echo "$(tput setaf 1)""Deleting \"${tempAURDirectory}\" directory...""$(tput sgr0)"
  
  rm -rf "${tempAURDirectory}"
fi

echo "$(tput setaf 2)""Creating \"${tempAURDirectory}\" directory...""$(tput sgr0)"

mkdir "${tempAURDirectory}"

cp build/archlinux/PKGBUILD "${tempAURDirectory}"
echo "$(tput setaf 3)""Changing version""$(tput sgr0)"

sed -i "s/VERSION_PLACEHOLDER/$version/" "${tempAURDirectory}/PKGBUILD"

echo
echo "$(tput setaf 6)"Make sure you have an account on: aur.archlinux.org "$(tput sgr0)"
read -rep "Proceed? (y/n)" -i "y" answer

if [[ "${answer}" != "y" ]]; then
  exit 1
fi

if ! [ -f "$HOME/.ssh/aur" ]; then
  echo "You don't have required ssh keys. Generating..."
  ssh-keygen -f ~/.ssh/aur
  echo "Make sure paste your public key from '~/.ssh/aur.pub' into your 'aur.archlinux.org' account"
  echo "$(tput setaf 5)"Generating ssh config file"$(tput sgr0)"

  if ! [ -f ~/.ssh/config ]; then
cat >"$HOME/.ssh/config" <<EOF
Host aur.archlinux.org
IdentityFile ~/.ssh/aur
User aur
EOF
fi
  
fi

cd "$tempAURDirectory" || exit 2

echo "git clone ssh://aur@aur.archlinux.org/secret-service.git"
git clone ssh://aur@aur.archlinux.org/secret-service.git

# Necessary???
# echo "git remote add origin: ssh://aur@aur.archlinux.org/secret-service.git"
# git remote add origin ssh://aur@aur.archlinux.org/secret-service.git

echo "checking PKGBUID..."
result=$(namcap PKGBUID)

if [[ "${result}" != "" ]]; then
  echo "$(tput setaf 1)"Malformed PKGBUILD"$(tput sgr0)"
  echo "${result}"
  exit 3
fi

# .gitinore with *
echo "*" > .gitinore

# create .SRCINFO
makepkg --printsrcinfo > .SRCINFO

# git add
git add -f PKGBUILD .SRCINFO
# git commit (version)
git commit -m "$version"
# Submit AUR package
git push
