# Triton Gateway-client
This is the triton-gateway-client that uses a scheduler.
Used the Text-to-Image feature of Stable Diffusion.   

## 1. Docker Start
### 1.1 Download
```
git clone https://github.com/ahr-i/triton-gateway-client.git
```

### 1.2 Setting
```
cd triton-gateway-client
vim setting/setting.go
```
Modify the contents of the file.   
```
package setting

/* ----- Server Setting ----- */
const ServerPort string = "6000" // Edit this

/* ----- Scheduler Server Setting ----- */
// If you are not using a scheduler, change the 'SchedulerActive' variable to false.
const SchedulerActive bool = false           // Edit this
const SchedulerUrl string = "localhost:8000" // Edit this

// If you are not using a Scheduler, please set the AgentURL.
const AgentURL string = "localhost:7000" // Edit this
```

### 1.3 Build
```
docker build -t triton-gateway .
```

### 1.4 Run
```
docker run -it --rm --name triton-gateway --network host triton-gateway
```
