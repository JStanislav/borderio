const fruits = [
    "watermelon",
    "pineapple",
    "kiwi",
    "mango",
    "coconut",
    "papaya",
    "raspberry",
    "peach",
    "passion fruit",
    "banana",
]
const animals = [
    "elephant",
    "crocodile",
    "octopus",
    "giraffe",
    "hippopotamus",
    "bat",
    "platypus",
    "flamingo",
    "hedgehog",
    "armadillo",
]

const entities = [
    ...animals, ...fruits
]

const funnyAdjectives = [
    "turbo",
    "cosmic",
    "glowing",
    "grumpy",
    "elastic",
    "explosive",
    "nostalgic",
    "diabolical",
    "quantum",
    "melancholic",
    "supersonic",
    "furry",
    "invisible",
    "rebellious",
    "vintage",
    "hypnotic",
    "fierce",
    "gelatinous",
    "existentialist",
    "musical",
]

export const generateRandomName = () => {
    const randomAdjective = funnyAdjectives[Math.floor(Math.random() * funnyAdjectives.length)];
    const randomEntity = entities[Math.floor(Math.random() * entities.length)];
    return `${randomAdjective} ${randomEntity}`;
}