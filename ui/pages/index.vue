<script>
export default {
  data() {
    return {
      newTask: "",
      newPriority: "High", // Default priority
      todos: [],
      statusMessage: "",
      priorityFilter: "", // Filter by priority
      searchQuery: "", // Search by task name
    };
  },
  computed: {
    incompleteTodos() {
      return this.todos.filter((todo) => todo.status !== "done");
    },
    completeTodos() {
      return this.todos.filter((todo) => todo.status === "done");
    },
    sortedFilteredIncompleteTodos() {
      return this.incompleteTodos
        .filter((todo) => {
          const matchesPriority = this.priorityFilter
            ? todo.priority === this.priorityFilter
            : true;
          const matchesSearch = todo.task.toLowerCase().includes(
            this.searchQuery.toLowerCase()
          );
          return matchesPriority && matchesSearch;
        })
        .sort((a, b) => {
          const priorityOrder = { High: 1, Medium: 2, Low: 3 };
          return priorityOrder[a.priority] - priorityOrder[b.priority];
        });
    },
  },
  mounted() {
    this.fetchTodos();
  },
  methods: {
    async fetchTodos() {
      try {
        const response = await fetch('/api/v1/todos', {mode: 'no-cors'});
        if (!response.ok) {
          throw new Error('Failed to fetch todos from server');
        }
        const data = await response.json();
        const allTodo = [...data.incomplete_tasks]

        this.todos = allTodo;

      } catch (error) {
        console.error(error);
        this.statusMessage = "Could not fetch tasks from the server.";
        // Fallback to localStorage if API call fails
        const storedTodos = localStorage.getItem("todos");
        if (storedTodos) {
          this.todos = JSON.parse(storedTodos);
        } else {
          this.todos = [];
        }
      }
    },
    saveTodos() {
      localStorage.setItem("todos", JSON.stringify(this.todos));
    },
    async addTodo() {
      if (this.newTask.trim() === "") return;

      const newTodo = {

        task: this.newTask,
        status: "created",
        priority: this.newPriority,

      };

      try {
        const response = await fetch('/api/v1/todos', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(newTodo),
        });
        if (!response.ok) {
          throw new Error('Failed to add new task');
        }
        const addedTodo = await response.json();
        this.todos.push(addedTodo);
        this.newTask = "";
        this.newPriority = "High"; // Reset to default
        this.statusMessage = "Task added successfully.";
        this.saveTodos();
      } catch (error) {
        console.error(error);
        this.statusMessage = "Could not add task to the server.";
      }
    },
    enableEdit(todo) {
      todo.isEditing = true;
    },
    async editTodo(todo) {
      todo.isEditing = false;

      try {
        const response = await fetch(`/api/v1/todos/${todo.id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(todo),
        });
        if (!response.ok) {
          throw new Error('Failed to edit task');
        }
        this.statusMessage = "Task edited successfully.";
        this.saveTodos();
      } catch (error) {
        console.error(error);
        this.statusMessage = "Could not update task on the server.";
      }
    },
    async updateStatus(todo) {
      todo.status = todo.status === "done" ? "created" : "done";

      try {
        const response = await fetch(`/api/v1/todos/${todo.id}`, {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ Status: todo.status }),
        });
        if (!response.ok) {
          throw new Error('Failed to update task status');
        }
        this.statusMessage = "Task status updated.";
        this.saveTodos();
      } catch (error) {
        console.error(error);
        this.statusMessage = "Could not update task status on the server.";
      }
    },
    async deleteTodo(id) {
      try {
        const response = await fetch(`/api/v1/todos/${id}`, {
          method: 'DELETE',
        });
        if (!response.ok) {
          throw new Error('Failed to delete task');
        }
        this.todos = this.todos.filter((todo) => todo.id !== id);
        this.statusMessage = "Task deleted successfully.";
        this.saveTodos();
      } catch (error) {
        console.error(error);
        this.statusMessage = "Could not delete task from the server.";
      }
    },
  },
};
</script>
<template>
  <div class="todo-main">
    <h1>TODO List</h1>
    <div v-if="statusMessage" class="status-message">{{ statusMessage }}</div>
    <div class="input-group">
      <input
        v-model="newTask"
        placeholder="Enter new task"
        @keyup.enter="addTodo"
      />
      <select v-model="newPriority">
        <option value="High">High Priority</option>
        <option value="Medium">Medium Priority</option>
        <option value="Low">Low Priority</option>
      </select>
      <button @click="addTodo">Add</button>
    </div>

    <div class="filter-group">
      <label for="priorityFilter">Filter by Priority:</label>
      <select v-model="priorityFilter" id="priorityFilter">
        <option value="">All</option>
        <option value="High">High</option>
        <option value="Medium">Medium</option>
        <option value="Low">Low</option>
      </select>
    </div>

    <div class="search-group">
      <input v-model="searchQuery" placeholder="Search by task name" />
    </div>

    <div v-if="incompleteTodos.length > 0">
      <h2>Incomplete Tasks</h2>
      <div
        v-for="todo in sortedFilteredIncompleteTodos"
        :key="todo.id"
        class="todo-item"
      >
        <input
          v-if="todo.isEditing"
          v-model="todo.task"
          class="edit-input"
          @blur="editTodo(todo)"
          @keyup.enter="editTodo(todo)"
        />
        <span
          v-else
          :class="{ 'done-task': todo.status === 'done' }"
          @click="enableEdit(todo)"
          >{{ todo.task }}</span
        >
        <div class="buttons">
          <button
            :class="{ done: todo.status === 'done' }"
            @click="updateStatus(todo)"
          >
            ✔️
          </button>
          <button class="delete-button" @click="deleteTodo(todo.id)">🗑️</button>
        </div>
      </div>
    </div>

    <div v-if="completeTodos.length > 0">
      <h2>Complete Tasks</h2>
      <div v-for="todo in completeTodos" :key="todo.id" class="todo-item">
        <span :class="{ 'done-task': todo.status === 'done' }">{{
          todo.task
        }}</span>
        <div class="buttons">
          <button
            :class="{ done: todo.status === 'done' }"
            @click="updateStatus(todo)"
          >
            ✔️
          </button>
          <button class="delete-button" @click="deleteTodo(todo.id)">🗑️</button>
        </div>
      </div>
    </div>

    <div v-else>
      <p>No tasks available.</p>
    </div>
  </div>
</template>

<style scoped>
.todo-main {
  max-width: 400px;
  margin: 20px auto;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: #fff;
}

.input-group {
  display: flex;
  margin-bottom: 20px;
}

input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-right: 10px;
}

button {
  padding: 10px;
  border: none;
  background-color: #333;
  color: #fff;
  border-radius: 4px;
  cursor: pointer;
}

.filter-group,
.search-group {
  margin-bottom: 20px;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.buttons button {
  background-color: #f0f0f0;
  color: #333;
  margin-left: 5px;
  border-radius: 4px;
  padding: 5px 10px;
  transition: background-color 0.3s, color 0.3s;
}

.buttons button.done {
  background-color: #333;
  color: #fff;
}

.status-message {
  margin-bottom: 20px;
  padding: 10px;
  background-color: #f0f0f0;
  border-radius: 4px;
}

.done-task {
  text-decoration: line-through;
}

.edit-input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
</style>
