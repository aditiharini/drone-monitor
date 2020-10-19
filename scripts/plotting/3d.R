library(tidyverse)
library(gridExtra)
library(plotly)

combined <- read.csv("/Users/aditisri/Dropbox (MIT)/Drone-Project/measurements/combined_traces/home-10-14-1/processed/stats/srv-1602679652.csv")

plot_ly(combined, x=~latitude, y=~longitude, z=~altitude, color=~uplinkBw)
