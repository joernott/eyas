#
# spec file for eyas
#
Name:           eyas
Version:        %{version}
Release:        %{release}
Summary:        Eyaml encryption server
License:        BSD
Group:          Sytem/Utilities
Vendor:         Ott-Consult UG
Packager:       Joern Ott
Url:            https://github.com/joernott/eyas
Source0:        eyas-%{version}.tar.gz
BuildArch:      x86_64

%description
A web UI for encrypting passwords for hiera using eyaml.

%prep
cd "$RPM_BUILD_DIR"
rm -rf *
tar -xzf "%{SOURCE0}"
STATUS=$?
if [ $STATUS -ne 0 ]; then
  exit $STATUS
fi
/usr/bin/chmod -Rf a+rX,u+w,g-w,o-w .

%build
cd "$RPM_BUILD_DIR/eyas-%{version}"
go get -u -v
go build -v

%install
install -Dpm 0755 %{name}-%{version}/eyas "%{buildroot}/usr/bin/eyas"

%files
%defattr(-,root,root,755)
%attr(755, root, root) /usr/bin/eyas

%changelog
* Mon Mar 06 2023 Joern Ott <joern.ott@ott-consult.de>
- Initial RPM version
