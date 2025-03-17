import { Box, Button, Container, Flex, Text } from "@chakra-ui/react";
import { useColorMode, useColorModeValue } from "./ui/color-mode";
import { IoMoon } from "react-icons/io5";
import { LuSun } from "react-icons/lu";


export default function Navbar() {
    const { colorMode, toggleColorMode } = useColorMode();

    return (
        <Container maxW={"1200px"}>
            <Box bg={useColorModeValue("gray.400", "gray.700")} px={4} my={4} borderRadius={4}>
                <Flex h={16} alignItems={"center"} justifyContent={"space-between"}>
                    <Flex
                        justifyContent={"center"} alignItems={"center"} gap={3} display={{ base: "none", sm: "flex" }}>
                            <img src="/react.png" alt="logo" width={50} height={50} />
                            <Text fontSize={"40"}>|</Text>
                            <img src="/golang.png" alt="logo" width={50} height={50} />
                            <Text fontSize={"40"}>|</Text>
                            <img src="/go.png" alt="logo" width={50} height={50} />
                    </Flex>

                    <Flex alignItems={"center"} gap={3}>
                        <Text fontSize={"lg"} fontWeight={500}>Tasks</Text>
                        <Button onClick={toggleColorMode}>
                            {colorMode === "light" ? <IoMoon/> : <LuSun size={20}/>}
                        </Button>
                    </Flex>
                </Flex>
            </Box>
        </Container>
    )


}