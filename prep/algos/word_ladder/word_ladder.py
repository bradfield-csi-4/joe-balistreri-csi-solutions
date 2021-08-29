from collections import defaultdict, deque
from itertools import product
import os

def build_graph(words):
    buckets = defaultdict(list)
    graph = defaultdict(set)

    for word in words:
        for i in range(len(word)):
            bucket = '{}_{}'.format(word[:i], word[i+1:])
            buckets[bucket].append(word)

    for bucket, mutual_neighbors in buckets.items():
        for word1, word2 in product(mutual_neighbors, repeat=2):
            if word1 != word2:
                graph[word1].add(word2)
                graph[word2].add(word1)

    return graph


def get_words(vocab_file):
    for line in open(vocab_file, 'r'):
        yield line[:-1]


def traverse(graph, starting_vertex):
    visited = set()
    queue = deque([[starting_vertex]])
    while queue:
        path = queue.popleft()
        vertex = path[-1]
        yield vertex, path
        for neighbor in graph[vertex] - visited:
            visited.add(neighbor)
            queue.append(path + [neighbor])


vocab_file = "/tmp/vocabulary.txt"
word_graph = build_graph(get_words(vocab_file))


if __name__ == '__main__':
    for vertex, path in traverse(word_graph, "FOOL"):
        if vertex == 'SAGE':
            print ' -> '.join(path)
