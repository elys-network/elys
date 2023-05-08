pkill -f "docker-compose .*elys.* logs" | true
pkill -f "/bin/bash.*create_logs.sh" | true
pkill -f "tail .*.log" | true
