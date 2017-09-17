package com.marshmallow.anwork.task;

import com.marshmallow.anwork.core.ProtobufSerializer;
import com.marshmallow.anwork.core.Serializable;
import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.task.protobuf.TaskManagerProtobuf;
import com.marshmallow.anwork.task.protobuf.TaskProtobuf;

import java.util.Arrays;
import java.util.PriorityQueue;

/**
 * This guy in a public interface for managing {@link Task} instances.
 *
 * <p>
 * Created Aug 29, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManager implements Serializable<TaskManagerProtobuf> {

  /**
   * This is the singleton {@link Serializer} for this class.
   */
  public static Serializer<TaskManager> SERIALIZER
      = new ProtobufSerializer<TaskManagerProtobuf, TaskManager>(() -> new TaskManager(),
                                                                 TaskManagerProtobuf.parser());

  private PriorityQueue<Task> tasks = new PriorityQueue<Task>();
  private Task currentTask;

  /**
   * Create a task from a name.
   *
   * @param name The name for the task
   * @param description The description for the task
   * @param priority The priority for the task
   * @throws IllegalArgumentException If this task already exists
   */
  public void createTask(String name,
                         String description,
                         int priority) throws IllegalArgumentException {
    if (findTask(name) != null) {
      throw new IllegalArgumentException("Task " + name + " already exists");
    }
    Task task = new Task(name, description, priority);
    tasks.add(task);
  }

  /**
   * Delete a task from a name.
   *
   * @param name The name for the task
   * @throws IllegalArgumentException If this task does not exist.
   */
  public void deleteTask(String name) throws IllegalArgumentException {
    Task task = findTask(name);
    if (task == null) {
      throw new IllegalArgumentException("Task " + name + " does not exist");
    }
    tasks.remove(task);
  }

  /**
   * Get the current task name.
   *
   * @return name The name of the current task, or <code>null</code> if there
   *     is no current task.
   */
  public String getCurrentTask() {
    return (currentTask == null ? null : currentTask.getName());
  }

  /**
   * Set the current task.
   *
   * @param name The name of the current task
   * @throws IllegalArgumentException If this task does not exist
   */
  public void setCurrentTask(String name) throws IllegalArgumentException {
    Task task = findTask(name);
    if (task == null) {
      throw new IllegalArgumentException("Task " + name + " does not exist");
    }
    currentTask = task;
  }

  /**
   * Get the state of a task.
   *
   * @param name The name of the task
   * @return A string representing the state of the task.
   * @throws IllegalArgumentException If the task does not exist
   */
  public TaskState getState(String name) throws IllegalArgumentException {
    Task task = findTask(name);
    if (task == null) {
      throw new IllegalArgumentException("Task " + name + " does not exist");
    }
    return task.getState();
  }

  /**
   * Set the state of a task.
   *
   * @param name The name of the current task
   * @param state The name of the state.
   * @throws IllegalArgumentException If this task does not exist or the state
   *     is invalid.
   */
  public void setState(String name, TaskState state) throws IllegalArgumentException {
    Task task = findTask(name);
    if (task == null) {
      throw new IllegalArgumentException("Task " + name + " does not exist");
    }
    task.setState(state);
  }

  /**
   * Get the tasks associated with this task manager.
   *
   * @return The tasks associated with this manager
   */
  public Task[] getTasks() {
    // TODO: this should be improved! We should write our own heap that can
    // return the elements to us in order.
    Task[] taskArray = tasks.toArray(new Task[0]);
    Arrays.sort(taskArray);
    return taskArray;
  }

  /**
   * Get a string representation of the tasks in this manager.
   * @return A string representation of the tasks in this manager.
   */
  @Override
  public String toString() {
    StringBuilder builder = new StringBuilder();
    Task[] taskArray = tasks.toArray(new Task[0]);
    for (Task task : taskArray) {
      if (task == currentTask) {
        builder.append('*');
      }
      builder.append(task);
    }

    return builder.toString();
  }

  private Task findTask(String name) {
    Task[] taskArray = tasks.toArray(new Task[0]);
    for (Task task : taskArray) {
      if (task.getName().equals(name)) {
        return task;
      }
    }
    return null;
  }

  @Override
  public TaskManagerProtobuf marshall() {
    TaskManagerProtobuf.Builder builder = TaskManagerProtobuf.newBuilder();
    Task[] taskArray = tasks.toArray(new Task[0]);
    for (Task task : taskArray) {
      builder.addTasks(task.marshall());
    }
    return builder.build();
  }

  @Override
  public void unmarshall(TaskManagerProtobuf protobuf) {
    for (int i = 0; i < protobuf.getTasksCount(); i++) {
      TaskProtobuf taskProtobuf = protobuf.getTasks(i);
      Task task = Task.FACTORY.makeBlankInstance();
      task.unmarshall(taskProtobuf);
      tasks.add(task);
    }
  }
}
