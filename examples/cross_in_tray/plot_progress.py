import json
import math
import os.path
import sys

from matplotlib import animation
import matplotlib.pyplot as plt
import numpy as np


PROGRESS_FILE_NAME = 'progress.json'
BOUNDS = [-10, 10]


def styblinski_tang(x, y):
    return -0.0001 * (
        np.abs(
            np.sin(x) * np.sin(y) * np.exp(
                np.abs(
                    100 - np.sqrt(x ** 2, y ** 2) / math.pi
                )
            )
        ) + 1
    ) ** 0.1


def main():

    # Check the progress file exists
    if not os.path.isfile(PROGRESS_FILE_NAME):
        print('{} does not exit, execute "go run main.go" first'.format(PROGRESS_FILE_NAME))
        sys.exit()

    # Load each line in the progress file
    steps = [json.loads(line) for line in open(PROGRESS_FILE_NAME)]

    # Extract the coordinates of each individual at each step
    coordinates = np.array([
        [
            indi['genome']
            for indi in step['pops'][0]['indis']
        ]
        for step in steps
    ])

    fig = plt.figure(figsize=(8, 8))
    ax = plt.axes(xlim=BOUNDS, ylim=BOUNDS)

    # Draw a heatmap of the function
    X, Y = np.meshgrid(
        np.linspace(BOUNDS[0], BOUNDS[1], 500),
        np.linspace(BOUNDS[0], BOUNDS[1], 500)
    )
    Z = styblinski_tang(X, Y)
    heatmap = ax.imshow(Z, extent=BOUNDS * 2)

    line, = ax.plot([], [], marker='o', color='orange', linestyle='None')

    def init():
        line.set_data([], [])
        return line,

    def animate(i):
        line.set_data(coordinates[i][:, 0], coordinates[i][:, 1])
        return line,

    anim = animation.FuncAnimation(fig, animate, init_func=init,
                               frames=len(coordinates), interval=100, blit=True)

    anim.save('progress.gif', writer='imagemagick', fps=10)

    plt.show()


if __name__ == '__main__':
    main()
