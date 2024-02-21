package setting

/* ----- Server Setting ----- */
const ServerPort string = "6000" // Edit this

/* ----- Scheduler Server Setting ----- */
// If you are not using a scheduler, change the 'SchedulerActive' variable to false.
const SchedulerActive bool = false           // Edit this
const SchedulerUrl string = "localhost:8000" // Edit this

// If you are not using a Scheduler, please set the AgentURL.
const AgentURL string = "localhost:7000" // Edit this
