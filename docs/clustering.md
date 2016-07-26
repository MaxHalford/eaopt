Clusters are a partitioning of individuals into smaller groups of similar individuals. Programmatically a cluster is a list of lists each containing individuals. Individuals inside each clusters are supposed to be similar. The similarity depends on a metric, for example it could be based on the fitness of the individuals.

!!! note "Note"
    In the litterature, clustering is also called **speciation**.

The purpose of a partinioning individuals is to apply genetic operators to similar individuals. In biological terms this encourages "incest" and maintains isolated species. For example in nature animals usually breed with local mates and don't breed with different animal species.

Using clustering/speciation with genetic algorithms became "popular" when they were first applied to the [optimization of neural network topologies](https://www.wikiwand.com/en/Neuroevolution_of_augmenting_topologies). By mixing two neural networks during crossover, the resulting neural networks were often useless because the inherited weights were not optimized for the new topology. This meant that newly generated neural networks were not performing well and would likely dissapear during selection. Thus speciation was introduced so that neural networks evolved in similar groups so that new neural networks wouldn't dissapear immediatly. Instead the similar neural networks would evolve between each other until they were good enough to mixed with the other neural networks.

With gago it's possible to use clustering on top of all the rest. For the while the only kind of clustering is fitness based. Later on it will be possible to provided a function to compare two individuals based on their genome. What happens is that a population of $n$ individuals is clustered into $k$ before applying an evolution model to each cluster. The $k$ clusters are then merged into a new population of $n$ individuals. This way each cluster didn't interact with the other clusters.

![clustering](img/clustering.png)
