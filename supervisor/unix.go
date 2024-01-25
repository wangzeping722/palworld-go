//go:build unix

package supervisor

import (
	"context"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Verthandii/palworld-go/config"
	"github.com/Verthandii/palworld-go/rcon"
)

type supervisor struct {
	c *rcon.Client
}

func New() (Supervisor, error) {
	cfg := config.CFG()
	c, err := rcon.New(cfg.Address, cfg.AdminPassword)
	if err != nil {
		return nil, err
	}

	return &supervisor{
		c: c,
	}, nil
}

func (s *supervisor) Start(ctx context.Context) {
	cfg := config.CFG()
	checkDuration := time.Duration(cfg.CheckInterval) * time.Second
	ticker := time.NewTicker(checkDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("退出 Supervisor")
			return
		case <-ticker.C:
			if !s.isAlive() {
				s.restart()
			}
			ticker.Reset(checkDuration)
		}
	}
}

func (s *supervisor) isAlive() bool {
	cfg := config.CFG()
	out, err := exec.Command("pgrep", "-f", cfg.ProcessName).Output()
	if err != nil {
		log.Println("Supervisor 健康检查失败", err)
		return false
	}

	// 检查输出结果，如果结果不为空，则至少存在一个进程
	output := strings.TrimSpace(string(out))
	return output != ""
}

func (s *supervisor) restart() {
	cfg := config.CFG()
	command := filepath.Join(cfg.GamePath, cfg.ProcessName)
	cmd := exec.Command(command)
	cmd.Dir = cfg.GamePath // 设置工作目录为游戏路径
	if err := cmd.Start(); err != nil {
		log.Printf("服务器重启失败【%v】\n", err)
	} else {
		log.Printf("服务器重启成功\n")
	}
}
