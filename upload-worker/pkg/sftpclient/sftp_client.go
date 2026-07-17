package sftpclient

import (
	"fmt"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func NewSFTPClient(cfg *configs.SftpClientConfig) (*sftp.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: cfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	sshClient, err := ssh.Dial("tcp", addr, sshConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH %w", err)
	}

	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("failed to start SFTP %w", err)
	}

	return client, nil
}
