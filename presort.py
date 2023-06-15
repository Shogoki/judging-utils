import os
import shutil

import markdown

def get_next_issue_number(prefix):
    """
    Get the next available issue number based on existing folders in the given directory.
    """
    existing_folders = [
        folder for folder in os.listdir() if folder.startswith(f"{prefix}-")
    ]
    if not existing_folders:
        return 1
    else:
        last_folder = max([int(f.replace(f"{prefix}-","")) for f in existing_folders] )
        return int(last_folder) + 1

def process_issue(issue_path, issue_type):
    """
    Process the given issue based on the selected issue type.
    """
    if issue_type == "i":
        shutil.move(issue_path, "invalid")
        print("Issue moved to 'invalid' folder.")
    elif issue_type == "m":
        next_number = get_next_issue_number("M")
        new_folder = f"M-{next_number:03d}"
        os.makedirs(new_folder)
        shutil.move(issue_path, new_folder)
        print(f"Issue moved to '{new_folder}' folder.")
    elif issue_type == "h":
        next_number = get_next_issue_number("H")
        new_folder = f"H-{next_number:03d}"
        os.makedirs(new_folder)
        shutil.move(issue_path, new_folder)
        print(f"Issue moved to '{new_folder}' folder.")
    elif issue_type == "d":
        existing_folders = [folder for folder in os.listdir() if os.path.isdir(folder) and (folder.startswith("M-") or folder.startswith("H-"))]
        existing_folders.sort()
        print("Existing folders:")
        for i, folder in enumerate(existing_folders, start=1):
            first_issue = os.listdir(folder)[0]
            with open(os.path.join(folder, first_issue), "r") as file:
                lines = file.readlines()
            title = lines[4].strip()
            print(f"{i}. {folder} - {title}")
        choice = input("Enter the number of the existing folder: ")
        if choice.isdigit() and int(choice) in range(1, len(existing_folders) + 1):
            existing_folder = existing_folders[int(choice) - 1]
            shutil.move(issue_path, existing_folder)
            print(f"Issue moved to '{existing_folder}' folder.")
        else:
            print("Invalid choice. Issue not moved.")
    else:
        print("Invalid choice. Issue not moved.")
        
def extract_summary(lines):
    """
    Extract the summary from the given list of lines.
    """
    summary_start = 7
    for i, line in enumerate(lines[summary_start:], start=summary_start):
        if line.startswith("##"):
            summary_end = i
            break
    else:
        summary_end = len(lines)
    summary = "".join(lines[summary_start:summary_end]).strip()
    return summary
         

def main():
    #thread.start_new_thread(start_server, ())
    print("Started Webserver in Background thread")
    # Create 'invalid' folder if it doesn't exist
    os.makedirs("invalid", exist_ok=True)

    # Get all markdown files in the current directory
    issues = [file for file in os.listdir() if file.endswith(".md")]

    if not issues:
        print("No issues found.")
        return
    issues.sort()
    # Iterate over the issues
    for issue in issues:
        with open(issue, "r") as file:
            lines = file.readlines()
        refresh = '<meta http-equiv="refresh" content="10" />' 
        html = markdown.markdown("\n".join(lines)) 
        with open("CURRENT.html", "w") as file:
            file.writelines([refresh,html])
        # Extract relevant information from the issue
        auditor = lines[0].strip().lstrip("#")
        severity = lines[2].strip()
        title = lines[4].strip()
        summary = extract_summary(lines)

        # Display issue information
        print("\n\nIssue:", issue)
        print("Title:", title)
        print("Summary:", summary)

        # Prompt for issue type and process the issue
        issue_type = input("Select the issue type ((i)nvalid/(m)edium/(h)igh/(d)uplicate/(s)kip): ")
        if issue_type == "s":
            continue
        elif issue_type == "q":
            exit()
        process_issue(issue, issue_type.lower())

        print()  # Print a blank line for separation

if __name__ == "__main__":
    main()

