# cpu.usage.total & cpu.usage.system

[<<BACK TO MAIN](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/README.md)

`cpu.usage.total` and `cpu.usage.system` are collected by SignalFx Smart Agent monitor type `docker-container-stats`.

Documentation: [Docker CPU Metrics](https://www.anshulpatel.in/post/linux_cpu_percentage/)

Further Reading: [What are CPU Jiffies?](https://www.anshulpatel.in/post/linux_cpu_percentage/)

Unlike host cpu utilization that can be calculated by subtracting the percentage used by the idle process from 100, container utilization is calculated by looking at how much time the container spends using the cpu out of the time spent by the entire system using the cpu.

When we chart the CPU Utilization % for a container, we need to apply a formula to the two metrics to obtain a percentage value. 

You will also find there may be more metric time series - `ts` - than actual containers at a given point in time, especially when you are looking at a longer time window. This is due to containers spinning up and down regularly, and time series is created for every container that existed - however it becomes inactive after some time when the container is no longer reporting.

![Docker CPU Utilization Chart](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/DockerCPUMetrics.png?raw=true)

_Click on the image to enlarge in a new tab_
