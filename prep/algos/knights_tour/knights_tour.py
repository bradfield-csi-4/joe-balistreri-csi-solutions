

def knights_tour(board_size):

    for x in range(board_size):
        for y in range(board_size):
            seen = set()
            stack = []
            success = tour(x ,y, seen, stack, board_size)
            if success:
                print(stack)
                return
            print("failed for ", x ,y)



def tour(x, y, seen, stack, board_size):
    seen.add((x, y))
    stack.append((x,y))

    if len(stack) == board_size * board_size:
        return True

    possible_moves = legal_moves(x, y, board_size) - seen

    for px, py in possible_moves:
        if tour(px, py, seen, stack, board_size):
            return True
        seen.remove((px, py))
        stack.pop()

    return False


def legal_moves(x, y, board_size):
    deltas = [
        (-1, -2), (-1, +2),
        (+1, -2), (+1, +2),
        (-2, -1), (+2, -1),
        (-2, +1), (+2, +1),
    ]
    moves = [(x + dx, y + dy) for dx, dy in deltas]
    return set([(x,y) for x, y in moves
        if x >= 0 and x <= board_size - 1 and y >= 0 and y <= board_size - 1])

if __name__ == '__main__':
    knights_tour(6)
