package com.marshmallow.anwork.app;

import java.io.File;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

import com.marshmallow.anwork.core.Cli;
import com.marshmallow.anwork.core.FilePersister;
import com.marshmallow.anwork.core.Persister;
import com.marshmallow.anwork.task.TaskManager;

/**
 * This is the main class for the anwork app.
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public class AnworkApp {

  public static void main(String[] args) {
    try {
      new AnworkApp().run(args);
    } catch (Exception e) {
      System.out.println("Error: " + e.getMessage());
    }
  }

  private static class Action {

    private static enum Type {
      CREATE,
      SHOW,
      ;
    }

    private final Type type;
    private final String taskName;

    public Action(Type type, String taskName) {
      this.type = type;
      this.taskName = taskName;
    }

    public Type getType() {
      return type;
    }

    public String getTaskName() {
      return taskName;
    }
  }

  private boolean debug = false;
  private String context = "default-context";
  private File persistenceRoot = new File(".");
  private List<Action> actions = new ArrayList<Action>();

  private void run(String[] args) throws Exception {
    Cli cli = makeCli();
    cli.runActions(args);

    TaskManager taskManager = loadTaskManager();

    runTaskAction(taskManager);

    saveTaskManager(taskManager);
  }

  private Cli makeCli() throws Exception {
    Cli cli = new Cli();
    cli.addAction("d",
                  "debug",
                  "Turn on extra debug printing",
                  null, // no argument name
                  (a) -> AnworkApp.this.debug = true);
    cli.addAction("c",
                  "context",
                  "Set the context in which this app runs",
                  "context-name", // no argument name
                  (a) -> AnworkApp.this.context = a);
    cli.addAction("o",
                  "output",
                  "Set the output file directory for persistant data",
                  "output-dir", // no argument name
                  (a) -> AnworkApp.this.persistenceRoot = new File(a));
    cli.addAction("t",
                  "task-create",
                  "Create a task",
                  "task-name",
                  (a) -> AnworkApp.this.actions.add(new Action(Action.Type.CREATE, a)));
    cli.addAction("s",
                  "show",
                  "Show the tasks that have been created",
                  null,
                  (a) -> AnworkApp.this.actions.add(new Action(Action.Type.SHOW, null)));
    return cli;
  }

  private TaskManager loadTaskManager() throws Exception {
    Persister<TaskManager> persister = new FilePersister<TaskManager>(persistenceRoot);
    if (!persister.contextExists(context)) {
      debugPrint("context " + context
                 + " does not exist at root " + persistenceRoot + "!"
                  + " creating new task manager");
      return new TaskManager();
    }

    Collection<TaskManager> loadeds = persister.load(context, TaskManager.serializer());
    if (loadeds.size() != 1) {
      throw new IllegalStateException("Persistence root " + persistenceRoot
                                      + " and context " + context
                                      + " contains 0 or more than 1 task manager.");
    }
    return loadeds.toArray(new TaskManager[0])[0];
  }

  private void runTaskAction(TaskManager taskManager) throws Exception {
    for (Action action : actions) {
      switch (action.getType()) {
      case CREATE:
        debugPrint("creating task " + action.getTaskName());
        taskManager.createTask(action.getTaskName(), "ummm", 1);
        break;
      case SHOW:
        System.out.println(taskManager.toString());
        break;
      default:
        throw new IllegalStateException("Unknown action type: " + action.getType());
      }
    }
  }

  private void saveTaskManager(TaskManager taskManager) throws Exception {
    Persister<TaskManager> persister = new FilePersister<TaskManager>(persistenceRoot);
    persister.save(context, TaskManager.serializer(), java.util.Collections.singleton(taskManager));
  }

  private void debugPrint(String string) {
    if (debug) {
      System.out.println("debug: " + string);
    }
  }
}
