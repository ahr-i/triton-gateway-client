# Triton Gateway-client
This is the triton-gateway-client that uses a scheduler.
Used the Text-to-Image feature of Stable Diffusion.   

## 1. Docker Start
### 1.1 Clone
```
git clone https://github.com/ahr-i/triton-gateway-client.git
```

### 1.2 build
```
cd triton-gateway-client
docker build -t triton-gateway .
```

### 1.3 setting
```
vim setting/setting.go
```
Modify the contents of the file.   
```
package setting

/* ----- Server Setting ----- */
const ServerPort string = "80" // Edit this

const ModelPath string = "./models/model_list.json"

/* ----- Triton Server Setting ----- */
const SchedulerUrl string = "localhost:8000" // Edit this

const batchSize int = 1
const Samples int = 1
const Steps int = 45
const GuidanceScale float64 = 7.5
```

### 1.4 Run
```
docker run -it --rm --name triton_gateway --network host triton-gateway
```
