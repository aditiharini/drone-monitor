library(tidyverse)
library(gridExtra)


combined <- read.csv("/Users/aditisri/Dropbox (MIT)/Drone-Project/measurements/combined_traces/home-10-8-1/processed/stats/srv-1602192535.csv")

p1 <- ggplot() +
  geom_point(data=combined, aes(x=time, y=uplinkBw, colour=factor(cellId))) + 
  xlab("time (sec)") + 
  ylab("bw (Mbps)") 

p2 <- ggplot() + 
  geom_point(data=combined, aes(x=time, y=rsrp, colour=factor(cellId))) + 
  xlab("time (sec)") + 
  ylab("rsrp (dBm)") 

grid.arrange(p1, p2)