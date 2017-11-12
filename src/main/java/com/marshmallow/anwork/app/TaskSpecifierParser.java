package com.marshmallow.anwork.app;

import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;

/**
 * This is a class that can parse special {@link Task} syntax passed on the command line and
 * convert it to an array of {@link Task} name's. These task names can then be used to perform
 * actions on {@link Task} instances using a {@link TaskManager}.
 *
 *<p>One or more {@link Task}'s can be specified on the command line via a couple of different ways.
 * <ol>
 * <li>The name of the {@link Task} (see {@link Task#getName()}) may be used. For example,
 *     "task-a."</li>
 * <li>The id of the {@link Task} (see {@link Task#getId()}) may be used. The ID of the task must
 *     start with an "at" sign: {@literal @}. For example, "@24."
 * </ol>
 *
 * <p>
 * Created Nov 11, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskSpecifierParser {

  /** This is the special character that comes before Task ID's. */
  public static final char SPECIAL_CHARACTER = '@';

  private final TaskManager manager;

  public TaskSpecifierParser(TaskManager manager) {
    this.manager = manager;
  }

  /**
   * Parse a command line string that refers to one or more {@link Task} names.
   *
   * <p>
   * Note! This method will *always* return a non-{@code null} {@link String} array of length
   * greater than 0. It will throw an {@link IllegalArgumentException} otherwise.
   * </p>
   *
   * @param string The command line argument
   * @return A non-{@code null} {@link String} array of length greater than 0 that contains the
   *     parsed names of the {@link Task}'s referred to by this command line argument
   * @throws IllegalArgumentException if the provided {@code string} cannot be parsed into one or
   *     more {@link Task} names
   */
  public String[] parse(String string) throws IllegalArgumentException {
    if (string.length() == 0) {
      throw new IllegalArgumentException("Cannot convert 0 length string to a task specifier");
    }

    // If the string starts with the special character, we are looking at a task ID.
    if (string.charAt(0) == SPECIAL_CHARACTER) {
      String id = string.substring(1);
      Task task = findTaskById(id);
      if (task == null) {
        throw new IllegalArgumentException("Unknown task for ID " + id);
      }
      return new String[] { task.getName() };
    }

    // Otherwise, we should be looking at a single task name.
    Task task = findTaskByName(string);
    if (task == null) {
      throw new IllegalArgumentException("Cannot find task with specifier " + string);
    }
    return new String[] { task.getName() };
  }

  private Task findTaskById(String id) {
    int numberId;
    try {
      numberId = Integer.parseInt(id);
    } catch (NumberFormatException exception) {
      return null;
    }

    for (Task task : manager.getTasks()) {
      if (task.getId() == numberId) {
        return task;
      }
    }

    return null;
  }

  private Task findTaskByName(String name) {
    for (Task task : manager.getTasks()) {
      if (task.getName().equals(name)) {
        return task;
      }
    }
    return null;
  }
}
