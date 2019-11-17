# KNN-Concurrency

![Alt Text](https://importq.files.wordpress.com/2017/11/knn_mov5.gif)

K Nearest Neighbour is a simple algorithm that stores all the available cases and classifies the new data or case based on a similarity measure. It is mostly used to classifies a data point based on how its neighbours are classified.
Let’s take below wine example. Two chemical components called Rutime and Myricetin. Consider a measurement of Rutine vs Myricetin level with two data points, Red and White wines. They have tested and where then fall on that graph based on how much Rutine and how much Myricetin chemical content present in the wines.

<img src="https://miro.medium.com/max/602/1*NDU4K7r-NIAt5lCN4-5EQQ.png"/>

‘k’ in KNN is a parameter that refers to the number of nearest neighbours to include in the majority of the voting process.
Suppose, if we add a new glass of wine in the dataset. We would like to know whether the new wine is red or white?

<img src="https://miro.medium.com/max/463/1*prEBTwv8V8BZiV-UbvXibQ.png"/>

So, we need to find out what the neighbours are in this case. Let’s say k = 5 and the new data point is classified by the majority of votes from its five neighbours and the new point would be classified as red since four out of five neighbours are red.
