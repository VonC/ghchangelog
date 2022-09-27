# ghchangelog

Convert GitHub changelog entry to markdown

## Example

    ghchangelog timezone
    
That will extract the article with '`timezone`' in its title (case insensitive), as listed at https://github.blog/changelog/

- if they are none, say so and exit
- if they are several, list the titles and exit
- it there is one (like "[Local timezones available on profiles](https://github.blog/changelog/2022-09-23-local-timezones-available-on-profiles/)"), convert it to markdown:

```markdown
> ## [Local timezones available on profiles](https://github.blog/changelog/2022-09-23-local-timezones-available-on-profiles) (Sep. 2022)
>
> You can now display your local timezone on your profile to give others an idea of when to expect responses to pull requests or issues from you.
> You can opt into this feature by navigating to Settings > Public Profile and checking `Display current local time`.
> You can also update this information directly from your profile by clicking 'Edit Profile' under your avatar.
>
> https://i0.wp.com/user-images.githubusercontent.com/4021812/191612405-01a07cf4-1280-4e79-9938-27415d0ed4b8.png?w=343&ssl=1 -- local timezone setting
>
> This will display your timezone in the left sidebar of your profile as well as your timezone's current deviation from UTC.
> When other users see your profile or user hovercard, they'll see your timezone as well as how many hours behind or ahead they are from your local time.
>
> https://i0.wp.com/user-images.githubusercontent.com/4021812/191612407-58d90e74-0cdb-4672-9686-8680f3355c18.png?w=535&ssl=1 -- local timezone display on profile
>
> [Learn more about personalizing your profile](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-github-profile/customizing-your-profile/personalizing-your-profile).
```

That would work with any portion of the title (again, case insensitive): 

    ghchangelog zones avail

That would also select the article titled "`Local timezones available on profiles`".  
(Assuming it is listed at https://github.blog/changelog, which was true at the time of writing)

No need for `"..."` (as in `ghchangelog "zones avail"`).  
`ghchangelog zones avail` is enough.

## installation

    go install github.com/VonC/ghchangelog@latest