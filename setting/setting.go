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
