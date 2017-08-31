package com.marshmallow.anwork.task;

import java.util.Arrays;
import java.util.PriorityQueue;

/**
 * This guy in a public interface for managing {@link Task} instances.
 *
 * @author Andrew
 * @date Aug 29, 2017
 */
public class TaskManager {

  private String context;
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
   * is no current task.
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
  public String getState(String name) throws IllegalArgumentException {
    Task task = findTask(name);
    if (task == null) {
      throw new IllegalArgumentException("Task " + name + " does not exist");
    }
    return task.getState().name().toLowerCase();
  }

  /**
   * Set the state of a task.
   *
   * @param name The name of the current task
   * @param state The name of the state.
   * @throws IllegalArgumentException If this task does not exist or the state
   * is invalid.
   */
  public void setState(String name, String state) throws IllegalArgumentException {
    Task task = findTask(name);
    if (task == null) {
      throw new IllegalArgumentException("Task " + name + " does not exist");
    }

    TaskState stateValue = TaskState.valueOf(state.toUpperCase());
    if (stateValue == null) {
      throw new IllegalArgumentException("State " + state + " is invalid. Here are the potential states: " + Arrays.toString(TaskState.values()));
    }

    task.setState(stateValue);
  }

  /**
   * Get the number of tasks that currently exist.
   *
   * @return The number of tasks that currently exist.
   */
  public int getTaskCount() {
    return tasks.size();
  }

  /**
   * Get a string representation of the tasks in this manager.
   * @return A string representation of the tasks in this manager.
   */
  @Override
  public String toString() {
    StringBuilder builder = new StringBuilder(context);
    builder.append(':');

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
}
