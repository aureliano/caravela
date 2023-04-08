package caravela

var ptBrMessages = map[int]string{
	100: "Baixando pacote de atualização.",
	101: "Baixando checksum.",
	200: "Verificando a existência de versão mais recente.",
	201: "Encontrada a versão %s",
	202: "Define %s como diretório de instalação.",
	203: "Descomprimindo arquivos.",
	204: "%d arquivos descomprimidos de %s.",
	205: "Apagando arquivos de instalação.",
	300: "Copiar %s para %s",
}

var enMessages = map[int]string{
	100: "Downloading update package.",
	101: "Downloading checksum file.",
	200: "Checking for the latest version.",
	201: "Found the version %s",
	202: "Sets %s as the installation directory.",
	203: "Decompressing files.",
	204: "%d decompressed files from %s.",
	205: "Deleting installation files.",
	300: "Copy %s to %s",
}
