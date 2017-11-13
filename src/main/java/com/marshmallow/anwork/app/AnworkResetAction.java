package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Action;
import com.marshmallow.anwork.app.cli.ArgumentValues;
import com.marshmallow.anwork.core.Persister;
import com.marshmallow.anwork.task.TaskManager;

import java.io.IOException;
import java.util.Scanner;

/**
 * This is a {@link Action} that completely deleted the persistence context on which this ANWORK
 * instance is running.
 *
 * <p>
 * Created Nov 12, 2017
 * </p>
 *
 * @author Andrew
 */
public class AnworkResetAction extends TaskManagerCliAction {
  @Override
  public boolean run(AnworkAppConfig config,
                     ArgumentValues flags,
                     ArgumentValues arguments,
                     TaskManager manager) {
    boolean force = flags.containsKey("f");
    if (force || doReallyActuallyDeleteEveryting()) {
      String context = config.getContext();
      Persister<TaskManager> persister = getPersister(config);
      try {
        if (!persister.exists(context)) {
          config.getDebugPrinter().accept("context " + context + " does not exist, "
                                          + " so there is nothing to delete!");
        } else {
          persister.clear(context);
        }
      } catch (IOException ioe) {
        System.out.println("Unable to reset: " + ioe.getMessage());
      }
    } else {
      config.getDebugPrinter().accept("Not erasing everything");
    }
    return false;
  }

  private boolean doReallyActuallyDeleteEveryting() {
    String answer;
    try (Scanner scanner = new Scanner(System.in)) {
      do {
        System.out.print("Are you sure you want to delete everything (y/n): ");
        answer = scanner.nextLine();
      } while (!(answer.equals("y") || answer.equals("n")));
    }
    return answer.equals("y");
  }
}
