import { Container, Stack } from '@chakra-ui/react'
import Navbar from './components/Navbar';
import TodoForm from './components/TodoForm';
import TodoList from './components/TodoList';

const BASE_URL = import.meta.env.MODE == 'development' ? 'http://127.0.0.1:8001' : 'http://127.0.0.1:8001';

function App() {

  return (
    <Stack h="100vh">
      <Navbar />
      <Container>
        <TodoForm />
        <TodoList />
      </Container>
    </Stack>
  );
}

export default App;
export { BASE_URL };