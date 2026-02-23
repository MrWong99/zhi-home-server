package main

import (
	"context"
	"sync"

	"github.com/MrWong99/zhi/pkg/zhiplugin/config"
)

// valueDefs contains all configuration values for the home server workspace.
var valueDefs = []ValueDef{
	// ── core ──────────────────────────────────────────────────────────────
	{
		Path: "core/timezone", Default: "Europe/Berlin",
		Section: "General", DisplayName: "Timezone",
		Description: "System timezone for all containers (TZ database name)",
		Type:        "string", Placeholder: "Europe/Berlin",
	},
	{
		Path: "core/domain", Default: "",
		Section: "General", DisplayName: "Base Domain",
		Description: "Base domain for service URLs (e.g., home.example.com)",
		Type:        "string", Placeholder: "home.example.com", Required: true,
	},
	{
		Path: "core/data-root", Default: "/srv/homeserver",
		Section: "Storage", DisplayName: "Data Root Path",
		Description: "Base directory for bind-mounted service data on the host",
		Type:        "string", Placeholder: "/srv/homeserver",
	},
	{
		Path: "core/compose-project-name", Default: "home-server",
		Section: "General", DisplayName: "Compose Project Name",
		Description: "Docker Compose project name (used for container/network naming)",
		Type:        "string",
	},

	// ── pihole ────────────────────────────────────────────────────────────
	{
		Path: "pihole/image-tag", Default: "latest",
		Section: "Image", DisplayName: "PiHole Image Tag",
		Description: "Docker image tag for pihole/pihole",
		Type:        "string",
	},
	{
		Path: "pihole/dns-port", Default: 53,
		Section: "Network", DisplayName: "DNS Port",
		Description: "Host port for DNS (UDP/TCP)",
		Type:        "int",
	},
	{
		Path: "pihole/web-port", Default: 8053,
		Section: "Network", DisplayName: "Web Admin Port",
		Description: "Host port for PiHole web admin interface",
		Type:        "int",
	},
	{
		Path: "pihole/admin-password", Default: "",
		Section: "Security", DisplayName: "Admin Password",
		Description: "Password for the PiHole web admin interface",
		Type:        "string", Password: true, Required: true,
	},
	{
		Path: "pihole/upstream-dns", Default: "1.1.1.1;8.8.8.8",
		Section: "DNS", DisplayName: "Upstream DNS Servers",
		Description: "Semicolon-separated upstream DNS servers",
		Type:        "string", Placeholder: "1.1.1.1;8.8.8.8",
	},
	{
		Path: "pihole/dnssec", Default: true,
		Section: "DNS", DisplayName: "Enable DNSSEC",
		Description: "Enable DNSSEC validation for DNS queries",
		Type:        "bool",
	},
	{
		Path: "pihole/custom-blocklists", Default: "",
		Section: "DNS", DisplayName: "Custom Blocklists",
		Description: "Comma-separated URLs for additional blocklists",
		Type:        "string", Placeholder: "https://example.com/blocklist.txt",
	},

	// ── plex ──────────────────────────────────────────────────────────────
	{
		Path: "plex/image-tag", Default: "latest",
		Section: "Image", DisplayName: "Plex Image Tag",
		Description: "Docker image tag for linuxserver/plex",
		Type:        "string",
	},
	{
		Path: "plex/web-port", Default: 32400,
		Section: "Network", DisplayName: "Web UI Port",
		Description: "Host port for Plex web interface",
		Type:        "int",
	},
	{
		Path: "plex/claim-token", Default: "",
		Section: "Account", DisplayName: "Claim Token",
		Description: "Plex claim token from https://plex.tv/claim (valid 4 minutes)",
		Type:        "string", Password: true, Placeholder: "claim-xxxxxxxxxxxxxxxxxxxx",
	},
	{
		Path: "plex/puid", Default: 1000,
		Section: "Permissions", DisplayName: "PUID",
		Description: "User ID for file ownership inside the container",
		Type:        "int",
	},
	{
		Path: "plex/pgid", Default: 1000,
		Section: "Permissions", DisplayName: "PGID",
		Description: "Group ID for file ownership inside the container",
		Type:        "int",
	},
	{
		Path: "plex/media-movies", Default: "/mnt/media/movies",
		Section: "Media Libraries", DisplayName: "Movies Path",
		Description: "Host path to movies library",
		Type:        "string", Placeholder: "/mnt/media/movies",
	},
	{
		Path: "plex/media-tv", Default: "/mnt/media/tv",
		Section: "Media Libraries", DisplayName: "TV Shows Path",
		Description: "Host path to TV shows library",
		Type:        "string", Placeholder: "/mnt/media/tv",
	},
	{
		Path: "plex/media-music", Default: "/mnt/media/music",
		Section: "Media Libraries", DisplayName: "Music Path",
		Description: "Host path to music library",
		Type:        "string", Placeholder: "/mnt/media/music",
	},
	{
		Path: "plex/hardware-transcoding", Default: false,
		Section: "Transcoding", DisplayName: "Hardware Transcoding",
		Description: "Enable hardware transcoding via /dev/dri (Intel Quick Sync / AMD VCE)",
		Type:        "bool",
	},

	// ── nextcloud ─────────────────────────────────────────────────────────
	{
		Path: "nextcloud/image-tag", Default: "latest",
		Section: "Image", DisplayName: "Nextcloud Image Tag",
		Description: "Docker image tag for nextcloud",
		Type:        "string",
	},
	{
		Path: "nextcloud/web-port", Default: 8080,
		Section: "Network", DisplayName: "Web Port",
		Description: "Host port for Nextcloud web interface",
		Type:        "int",
	},
	{
		Path: "nextcloud/admin-user", Default: "admin",
		Section: "Admin Account", DisplayName: "Admin Username",
		Description: "Nextcloud admin username (set during first run only)",
		Type:        "string",
	},
	{
		Path: "nextcloud/admin-password", Default: "",
		Section: "Admin Account", DisplayName: "Admin Password",
		Description: "Nextcloud admin password (set during first run only)",
		Type:        "string", Password: true, Required: true,
	},
	{
		Path: "nextcloud/trusted-domains", Default: "localhost",
		Section: "Security", DisplayName: "Trusted Domains",
		Description: "Space-separated list of trusted domains for Nextcloud",
		Type:        "string", Placeholder: "localhost cloud.home.example.com",
	},
	{
		Path: "nextcloud/max-upload-size", Default: "16G",
		Section: "Uploads", DisplayName: "Max Upload Size",
		Description: "Maximum file upload size (e.g., 512M, 1G, 16G)",
		Type:        "string",
	},
	{
		Path: "nextcloud/redis-file-locking", Default: true,
		Section: "Performance", DisplayName: "Redis File Locking",
		Description: "Use Redis for transactional file locking (recommended when Redis is available)",
		Type:        "bool",
	},
	{
		Path: "nextcloud/smtp-host", Default: "",
		Section: "Email (SMTP)", DisplayName: "SMTP Host",
		Description: "SMTP server hostname for sending emails",
		Type:        "string", Placeholder: "smtp.example.com",
	},
	{
		Path: "nextcloud/smtp-port", Default: 587,
		Section: "Email (SMTP)", DisplayName: "SMTP Port",
		Description: "SMTP server port",
		Type:        "int",
	},
	{
		Path: "nextcloud/smtp-user", Default: "",
		Section: "Email (SMTP)", DisplayName: "SMTP Username",
		Description: "SMTP authentication username",
		Type:        "string",
	},
	{
		Path: "nextcloud/smtp-password", Default: "",
		Section: "Email (SMTP)", DisplayName: "SMTP Password",
		Description: "SMTP authentication password",
		Type:        "string", Password: true,
	},

	// ── mariadb ───────────────────────────────────────────────────────────
	{
		Path: "mariadb/image-tag", Default: "11",
		Section: "Image", DisplayName: "MariaDB Image Tag",
		Description: "Docker image tag for mariadb",
		Type:        "string",
	},
	{
		Path: "mariadb/root-password", Default: "",
		Section: "Security", DisplayName: "Root Password",
		Description: "MariaDB root password (set during first run only)",
		Type:        "string", Password: true, Required: true,
	},
	{
		Path: "mariadb/nextcloud-db", Default: "nextcloud",
		Section: "Nextcloud Database", DisplayName: "Database Name",
		Description: "Database name for Nextcloud",
		Type:        "string",
	},
	{
		Path: "mariadb/nextcloud-user", Default: "nextcloud",
		Section: "Nextcloud Database", DisplayName: "Database User",
		Description: "Database user for Nextcloud",
		Type:        "string",
	},
	{
		Path: "mariadb/nextcloud-password", Default: "",
		Section: "Nextcloud Database", DisplayName: "Database Password",
		Description: "Database password for the Nextcloud user",
		Type:        "string", Password: true, Required: true,
	},
	{
		Path: "mariadb/enable-binlog", Default: false,
		Section: "Replication", DisplayName: "Enable Binary Logging",
		Description: "Enable binary logging for replication. Disabled by default to save disk space on home servers.",
		Type:        "bool",
	},
	{
		Path: "mariadb/innodb-buffer-pool-size", Default: "256M",
		Section: "Tuning", DisplayName: "InnoDB Buffer Pool Size",
		Description: "InnoDB buffer pool size (e.g., 256M, 1G)",
		Type:        "string", Placeholder: "256M",
	},

	// ── redis ─────────────────────────────────────────────────────────────
	{
		Path: "redis/image-tag", Default: "8-alpine",
		Section: "Image", DisplayName: "Redis Image Tag",
		Description: "Docker image tag for redis",
		Type:        "string",
	},
	{
		Path: "redis/maxmemory", Default: "128mb",
		Section: "Memory", DisplayName: "Max Memory",
		Description: "Maximum memory Redis can use (e.g., 128mb, 256mb)",
		Type:        "string", Placeholder: "128mb",
	},
	{
		Path: "redis/maxmemory-policy", Default: "allkeys-lru",
		Section: "Memory", DisplayName: "Eviction Policy",
		Description: "How Redis evicts keys when maxmemory is reached",
		Type:        "string",
		SelectFrom:  []string{"allkeys-lru", "volatile-lru", "allkeys-lfu", "volatile-lfu", "noeviction"},
	},

	// ── nginx-proxy-manager ───────────────────────────────────────────────
	{
		Path: "nginx-proxy-manager/image-tag", Default: "latest",
		Section: "Image", DisplayName: "NPM Image Tag",
		Description: "Docker image tag for jc21/nginx-proxy-manager",
		Type:        "string",
	},
	{
		Path: "nginx-proxy-manager/http-port", Default: 80,
		Section: "Ports", DisplayName: "HTTP Port",
		Description: "Host port for HTTP traffic",
		Type:        "int",
	},
	{
		Path: "nginx-proxy-manager/https-port", Default: 443,
		Section: "Ports", DisplayName: "HTTPS Port",
		Description: "Host port for HTTPS traffic",
		Type:        "int",
	},
	{
		Path: "nginx-proxy-manager/admin-port", Default: 81,
		Section: "Ports", DisplayName: "Admin UI Port",
		Description: "Host port for NPM admin web interface",
		Type:        "int",
	},
	{
		Path: "nginx-proxy-manager/letsencrypt-email", Default: "",
		Section: "SSL", DisplayName: "Let's Encrypt Email",
		Description: "Email for Let's Encrypt certificate notifications",
		Type:        "string", Placeholder: "admin@example.com",
	},
}

// homeserverPlugin implements config.Plugin.
type homeserverPlugin struct {
	mu     sync.RWMutex
	paths  []string
	values map[string]*config.Value
}

func newHomeserverPlugin() *homeserverPlugin {
	p := &homeserverPlugin{
		values: make(map[string]*config.Value, len(valueDefs)),
		paths:  make([]string, 0, len(valueDefs)),
	}
	for _, v := range valueDefs {
		p.paths = append(p.paths, v.Path)
		p.values[v.Path] = v.ToValue()
	}
	return p
}

func (p *homeserverPlugin) List(_ context.Context) ([]string, error) {
	return p.paths, nil
}

func (p *homeserverPlugin) Get(_ context.Context, path string) (config.Value, bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	v, ok := p.values[path]
	if !ok {
		return config.Value{}, false, nil
	}
	return *v, true, nil
}

func (p *homeserverPlugin) Set(_ context.Context, path string, v config.Value) error {
	if err := config.ValidatePath(path); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.values[path] = &v
	return nil
}

func (p *homeserverPlugin) Validate(_ context.Context, path string, tree config.TreeReader) ([]config.ValidationResult, error) {
	if fn, ok := validators[path]; ok {
		v, found := tree.Get(path)
		if !found {
			return nil, nil
		}
		return fn(v, tree)
	}
	return nil, nil
}
