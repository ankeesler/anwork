package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliAction;
import com.marshmallow.anwork.app.cli.CliList;
import com.marshmallow.anwork.core.FilePersister;
import com.marshmallow.anwork.core.Persister;
import com.marshmallow.anwork.task.TaskManager;

import java.io.File;
import java.util.Collection;

/**
 * This is the main class for the anwork app.
 *
 * @author Andrew
 * Created Sep 9, 2017
 */
public class AnworkApp {

  /**
   * ANWORK main method.
   *
   * @param args Command line argument
   */
  public static void main(String[] args) {
    try {
      new AnworkApp().run(args);
    } catch (Exception e) {
      System.out.println("Error: " + e.getMessage());
    }
  }

  // Global CLI flags
  private boolean debug = false;
  private String context = "default-context";
  private File persistenceRoot = new File(".");

  private void run(String[] args) throws Exception {
    makeCli().parse(args);
  }

  private void debugPrint(String string) {
    if (debug) {
      System.out.println("debug: " + string);
    }
  }

  /*
   * Section - CLI Creation
   */

  private Cli makeCli() throws Exception {
    Cli cli = new Cli("anwork", "ANWORK CLI commands");
    CliList root = cli.getRoot();
    makeRootFlags(root);
    makeTaskCommands(root);
    return cli;
  }

  private void makeRootFlags(CliList root) {
    root.addLongFlag("d",
                     "debug",
                     "Turn on debug printing",
      (p) -> AnworkApp.this.debug = true);
    root.addLongFlagWithParameter("c",
                                  "context",
                                  "Set the persistence context",
                                  "name",
      (p) -> AnworkApp.this.context = p[0]);
    root.addLongFlagWithParameter("o",
                                  "output",
                                  "Set persistence output directory",
                                  "directory",
      (p) -> AnworkApp.this.persistenceRoot = new File(p[0]));
  }

  private void makeTaskCommands(CliList root) {
    CliList taskCommandList = root.addList("task", "Task commands...");

    CliAction createAction = new CliAction() {
      @Override
      public void run(String[] args) {
        try {
          TaskManager manager = loadTaskManager();
          manager.createTask(args[0], args[1], Integer.parseInt(args[2]));
          saveTaskManager(manager);
        } catch (Exception e) {
          System.out.println("Could not create task: " + e.getMessage());
        }
      }
    };
    taskCommandList.addCommand("create", "Create a task", createAction);
  }

  /*
   * Section - Task Manager Management
   */

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

  private void saveTaskManager(TaskManager taskManager) throws Exception {
    Persister<TaskManager> persister = new FilePersister<TaskManager>(persistenceRoot);
    persister.save(context, TaskManager.serializer(), java.util.Collections.singleton(taskManager));
  }
}
