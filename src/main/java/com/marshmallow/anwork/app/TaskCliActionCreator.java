package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Action;
import com.marshmallow.anwork.app.cli.ActionCreator;
import com.marshmallow.anwork.app.cli.ArgumentType;
import com.marshmallow.anwork.app.cli.ArgumentValues;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskState;

/**
 * This is a {@link ActionCreator} for the task commands in the ANWORK app.
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskCliActionCreator implements ActionCreator {

  private static class SetStateCliAction extends TaskManagerCliAction {

    private TaskState taskState;

    public SetStateCliAction(TaskState taskState) {
      super();
      this.taskState = taskState;
    }

    @Override
    public boolean run(AnworkAppConfig config,
                       ArgumentValues flags,
                       ArgumentValues arguments,
                       TaskManager manager) {
      String taskName = getTaskSpecifierArgument(manager, arguments);
      manager.setState(taskName, taskState);
      return true;
    }
  }

  private static class ShowCliAction extends TaskManagerCliAction {

    public static final ShowCliAction INSTANCE = new ShowCliAction();

    @Override
    public boolean run(AnworkAppConfig config,
                       ArgumentValues flags,
                       ArgumentValues arguments,
                       TaskManager manager) {
      boolean showShort = flags.containsKey("s");
      printTasksForState(TaskState.RUNNING, manager, showShort);
      printTasksForState(TaskState.BLOCKED, manager, showShort);
      printTasksForState(TaskState.WAITING, manager, showShort);
      printTasksForState(TaskState.FINISHED, manager, showShort);
      return true;
    }

    private void printTasksForState(TaskState state, TaskManager manager, boolean showShort) {
      System.out.println(state.name() + " tasks:");
      for (Task task : manager.getTasks()) {
        if (task.getState().equals(state)) {
          String stuff = (showShort
                          ? AnworkAppUtilities.makeTaskShortString(task, manager, "  ")
                          : AnworkAppUtilities.makeTaskLongString(task, manager, "  "));
          System.out.println(stuff);
        }
      }
    }
  }

  @Override
  public Action createAction(String commandName) {
    switch (commandName) {
      case "create":
        return new TaskManagerCliAction() {
          @Override
          public boolean run(AnworkAppConfig config,
                             ArgumentValues flags,
                             ArgumentValues arguments,
                             TaskManager manager) {
            String name = arguments.getValue("task-name", ArgumentType.STRING);
            String description = (flags.containsKey("e")
                                  ? flags.getValue("e", ArgumentType.STRING)
                                  : "");
            int priority = (int)(flags.containsKey("p")
                                 ? flags.getValue("p", ArgumentType.NUMBER)
                                 : Task.DEFAULT_PRIORITY);
            if (name.charAt(0) == TaskSpecifierParser.SPECIAL_CHARACTER) {
              String message = String.format("Task name cannot begin with special character '%c'",
                                             TaskSpecifierParser.SPECIAL_CHARACTER);
              throw new IllegalStateException(message);
            }
            manager.createTask(name, description, priority);
            config.getDebugPrinter().accept("created task '" + name + "'");
            return true;
          }
        };
      case "set-waiting":
        return new SetStateCliAction(TaskState.WAITING);
      case "set-blocked":
        return new SetStateCliAction(TaskState.BLOCKED);
      case "set-running":
        return new SetStateCliAction(TaskState.RUNNING);
      case "set-finished":
        return new SetStateCliAction(TaskState.FINISHED);
      case "set-priority":
        return new TaskManagerCliAction() {
          @Override
          public boolean run(AnworkAppConfig config,
                             ArgumentValues flags,
                             ArgumentValues arguments,
                             TaskManager manager) {
            String name = getTaskSpecifierArgument(manager, arguments);
            Long priority = arguments.getValue("priority", ArgumentType.NUMBER);
            manager.setPriority(name, priority.intValue());
            return true;
          }
        };
      case "delete":
        return new TaskManagerCliAction() {
          @Override
          public boolean run(AnworkAppConfig config,
                             ArgumentValues flags,
                             ArgumentValues arguments,
                             TaskManager manager) {
            String name = getTaskSpecifierArgument(manager, arguments);
            manager.deleteTask(name);
            config.getDebugPrinter().accept("deleted task '" + name + "'");
            return true;
          }
        };
      case "delete-all":
        return new TaskManagerCliAction() {
          @Override
          public boolean run(AnworkAppConfig config,
                          ArgumentValues flags,
                          ArgumentValues arguments,
                          TaskManager manager) {
            for (Task task : manager.getTasks()) {
              manager.deleteTask(task.getName());
            }
            config.getDebugPrinter().accept("deleted all tasks");
            return true;
          }
        };
      case "show":
        return ShowCliAction.INSTANCE;
      case "note":
        return new TaskManagerCliAction() {
          @Override
          public boolean run(AnworkAppConfig config,
                             ArgumentValues flags,
                             ArgumentValues arguments,
                             TaskManager manager) {
            String name = getTaskSpecifierArgument(manager, arguments);
            String note = arguments.getValue("note", ArgumentType.STRING);
            manager.addNote(name, note);
            return true;
          }
        };
      default:
        return null; // error!
    }
  }
}
