library(tidyverse)
library(gridExtra)

latency <- read.csv("/Users/aditisri/Dropbox (MIT)/Drone-Project/measurements/saturatr_traces/control/processed/stats/latency-1599793905.csv")
throughput <- read.csv("/Users/aditisri/Dropbox (MIT)/Drone-Project/measurements/saturatr_traces/control/processed/stats/throughput-1599793905.csv")

p1 <- ggplot() + 
  geom_point(data=latency, aes(x=time, y=latency.ms.,), colour="purple") + 
  xlab("time (ms)") + 
  ylab("latency (ms)")
p2 <- ggplot() + 
  geom_point(data=throughput, aes(x=time, y=throughput.Mbps.), colour="purple") + 
  xlab("time (ms)") + 
  ylab("thoughput (Mbps)")

grid.arrange(p1, p2, nrow=1)