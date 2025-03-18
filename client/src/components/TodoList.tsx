import { Flex, Spinner, Stack, Text } from "@chakra-ui/react";
import TodoItem from "./TodoItem";
import { useQuery } from "@tanstack/react-query";
import { BASE_URL } from "../App";

export type Todo = {
    id: number;
    title: string;
    description: string;
    completed: boolean;
};

const TodoList = () => {
    const { data: todos, isLoading } = useQuery<Todo[]>({
        queryKey: ["todos"],
        queryFn: async () => {
            try {
                const res = await fetch(BASE_URL + "/api/todos");
                const jsonResponse = await res.json();
                const data = jsonResponse.success ? jsonResponse.data : [];

                if (!res.ok || !jsonResponse.success) {
                    throw new Error(jsonResponse.error || "Something went wrong");
                }

                return data || [];
            } catch (error) {
                console.log(error);
            }
        }
    });

    return (
        <>
            <Text
                fontSize={"4xl"}
                textTransform={"uppercase"}
                fontWeight={"bold"}
                textAlign={"center"}
                my={2}
                bgGradient='linear(to-l, #0b85f8, #00ffff)'
                bgClip='text'
            >
                All Tasks
            </Text>
            {isLoading && (
                <Flex justifyContent={"center"} my={4}>
                    <Spinner size={"xl"} />
                </Flex>
            )}
            {!isLoading && todos?.length === 0 && (
                <Stack alignItems={"center"} gap='3'>
                    <Text fontSize={"xl"} textAlign={"center"} color={"gray.500"}>
                        All tasks completed! 🤞
                    </Text>
                    <img src='/go.png' alt='Go logo' width={70} height={70} />
                </Stack>
            )}
            <Stack gap={3}>
                {todos?.map((todo) => (
                    <TodoItem key={todo.id} todo={todo} />
                ))}
            </Stack>
        </>
    );
};
export default TodoList;