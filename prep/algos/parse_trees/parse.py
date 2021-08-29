import operator

OPERATORS = {
    '+': operator.add,
    '-': operator.sub,
    '*': operator.mul,
    '/': operator.truediv
}
LEFT_PAREN = '('
RIGHT_PAREN = ')'

def build_parse_tree(expression):
    tree = {}
    stack = [tree]
    node = tree
    for token in expression:
        if token == LEFT_PAREN:
            node['left'] = {}
            stack.append(node)
            node = node['left']
        elif token == RIGHT_PAREN:
            node = stack.pop()
        elif token in OPERATORS:
            node['val'] = token
            node['right'] = {}
            stack.append(node)
            node = node['right']
        else:
            node['val'] = int(token)
            parent = stack.pop()
            node = parent
    return tree

def construct_expression(parse_tree);
    if parse_tree is None:
        return ''

    left = construct_expression(parse_tree.get('left'))
    right = construct_expression(parse_tree.get('right'))
    val = parse_tree['val']

    if left and right:
        return '({}{}{})'.format(left, val, right)

    return val


def evaluate(tree):
    if not tree['left'] and not tree['right']:
        return tree['val']
    operate = OPERATORS[tree['val']]
    return operate(evaluate(tree['left']), evaluate(tree['right']))


def preorder(node):
    if node:
        print(node['val'])
        preorder(node.get('left'))
        preorder(node.get('right'))

def postorder(node):
    if node:
        postorder(node.get('left'))
        postorder(node.get('right'))
        print(node['val'])

def inorder(node):
    if node:
        inorder(node.get('left'))
        print(node['val'])
        inorder(node.get('right'))
