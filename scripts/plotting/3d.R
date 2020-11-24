install.packages("ggplot2", repos="http://cran.us.r-project.org")
install.packages("gridExtra", repos="http://cran.us.r-project.org")
install.packages("plotly", repos="http://cran.us.r-project.org")

library(ggplot2)
library(gridExtra)
library(plotly)

args = commandArgs(trailingOnly=TRUE)
combined <- read.csv("Drone-Project/measurements/combined_traces/home-10-14-1/stats/srv-1602679652.csv")


plot_ly(combined, x=~latitude, y=~longitude, z=~altitude, color=~uplinkBw)
